package line

import (
	"net/http"
	"strings"

	"github.com/albert-einshutoin/mockport/internal/adapter/httpx"
	"github.com/albert-einshutoin/mockport/internal/state"
)

func (r *routes) writeValidateRichMenu(w http.ResponseWriter, req *http.Request) {
	payload, err := decodePayload(req)
	if err != nil {
		writeDecodeError(w, err)
		return
	}
	if _, ok := payload["name"].(string); !ok {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"message": "The request body has 1 error(s)",
			"details": []map[string]string{{
				"message":  "must be specified",
				"property": "name",
			}},
		})
		return
	}
	writeEmptyJSON(w, http.StatusOK)
}

func (r *routes) writeCreateRichMenu(w http.ResponseWriter, req *http.Request) {
	payload, err := decodePayload(req)
	if err != nil {
		writeDecodeError(w, err)
		return
	}
	if _, ok := payload["name"].(string); !ok {
		writeLINEError(w, http.StatusBadRequest, "The request body has 1 error(s)")
		return
	}
	resource, err := r.store.Create("line", "rich_menu", payload)
	if err != nil {
		writeLINEError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"richMenuId": resource.ID})
}

func (r *routes) writeRichMenuList(w http.ResponseWriter) {
	menus := []map[string]any{}
	for _, resource := range r.store.List("line", "rich_menu") {
		menus = append(menus, richMenuResponse(resource))
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"richmenus": menus})
}

func (r *routes) writeGetRichMenu(w http.ResponseWriter, id string) {
	resource, ok := r.store.Get("line", "rich_menu", id)
	if !ok {
		writeLINEError(w, http.StatusNotFound, "Not found")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, richMenuResponse(resource))
}

func (r *routes) writeDeleteRichMenu(w http.ResponseWriter, id string) {
	if !r.deleteRichMenu(id) {
		writeLINEError(w, http.StatusNotFound, "Not found")
		return
	}
	writeEmptyJSON(w, http.StatusOK)
}

func (r *routes) writeUploadRichMenuImage(w http.ResponseWriter, path string) {
	id := strings.TrimSuffix(strings.TrimPrefix(path, "/v2/bot/richmenu/"), "/content")
	if _, ok := r.store.Get("line", "rich_menu", id); !ok {
		writeLINEError(w, http.StatusNotFound, "Not found")
		return
	}
	if _, err := r.store.Update("line", "rich_menu", id, map[string]any{"hasImage": true}); err != nil {
		writeLINEError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeEmptyJSON(w, http.StatusOK)
}

func (r *routes) writeDownloadRichMenuImage(w http.ResponseWriter, path string) {
	id := strings.TrimSuffix(strings.TrimPrefix(path, "/v2/bot/richmenu/"), "/content")
	resource, ok := r.store.Get("line", "rich_menu", id)
	if !ok || resource.Data["hasImage"] != true {
		writeLINEError(w, http.StatusNotFound, "Not found")
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("mockport rich menu image"))
}

func (r *routes) writeSetDefaultRichMenu(w http.ResponseWriter, id string) {
	if status, message := r.setDefaultRichMenuChecked(id); status != http.StatusOK {
		writeLINEError(w, status, message)
		return
	}
	writeEmptyJSON(w, http.StatusOK)
}

func (r *routes) writeGetDefaultRichMenu(w http.ResponseWriter) {
	id, ok := r.currentDefaultRichMenu()
	if !ok {
		writeLINEError(w, http.StatusNotFound, "no default richmenu")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"richMenuId": id})
}

func (r *routes) writeLinkRichMenuToUser(w http.ResponseWriter, path string) {
	parts := strings.Split(strings.TrimPrefix(path, "/v2/bot/user/"), "/richmenu/")
	if len(parts) != 2 {
		writeLINEError(w, http.StatusBadRequest, "The value for the 'userId' parameter is invalid")
		return
	}
	if _, ok := r.store.Get("line", "rich_menu", parts[1]); !ok {
		writeLINEError(w, http.StatusNotFound, "Not found")
		return
	}
	_, err := r.store.Create("line", "user_rich_menu", map[string]any{"userId": parts[0], "richMenuId": parts[1]})
	if err != nil {
		writeLINEError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeEmptyJSON(w, http.StatusOK)
}

func (r *routes) writeGetUserRichMenu(w http.ResponseWriter, path string) {
	userID := strings.TrimSuffix(strings.TrimPrefix(path, "/v2/bot/user/"), "/richmenu")
	for _, resource := range r.store.List("line", "user_rich_menu") {
		if resource.Data["userId"] == userID {
			httpx.WriteJSON(w, http.StatusOK, map[string]any{"richMenuId": resource.Data["richMenuId"]})
			return
		}
	}
	writeLINEError(w, http.StatusNotFound, "the user has no richmenu")
}

func (r *routes) writeUnlinkRichMenuFromUser(w http.ResponseWriter, path string) {
	userID := strings.TrimSuffix(strings.TrimPrefix(path, "/v2/bot/user/"), "/richmenu")
	for _, resource := range r.store.List("line", "user_rich_menu") {
		if resource.Data["userId"] == userID {
			r.store.Delete("line", "user_rich_menu", resource.ID)
		}
	}
	writeEmptyJSON(w, http.StatusOK)
}

func (r *routes) writeCreateRichMenuAlias(w http.ResponseWriter, req *http.Request) {
	payload, err := decodePayload(req)
	if err != nil {
		writeDecodeError(w, err)
		return
	}
	aliasID, _ := payload["richMenuAliasId"].(string)
	richMenuID, _ := payload["richMenuId"].(string)
	if aliasID == "" || richMenuID == "" {
		writeLINEError(w, http.StatusBadRequest, "The request body has 1 error(s)")
		return
	}
	if _, ok := r.store.Get("line", "rich_menu", richMenuID); !ok {
		writeLINEError(w, http.StatusBadRequest, "richmenu not found")
		return
	}
	if _, ok := r.findAlias(aliasID); ok {
		writeLINEError(w, http.StatusBadRequest, "conflict richmenu alias id")
		return
	}
	if _, err := r.store.Create("line", "rich_menu_alias", map[string]any{"richMenuAliasId": aliasID, "richMenuId": richMenuID}); err != nil {
		writeLINEError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeEmptyJSON(w, http.StatusOK)
}

func (r *routes) writeUpdateRichMenuAlias(w http.ResponseWriter, req *http.Request, aliasID string) {
	resource, ok := r.findAlias(aliasID)
	if !ok {
		writeLINEError(w, http.StatusNotFound, "richmenu alias not found")
		return
	}
	payload, err := decodePayload(req)
	if err != nil {
		writeDecodeError(w, err)
		return
	}
	richMenuID, _ := payload["richMenuId"].(string)
	if _, ok := r.store.Get("line", "rich_menu", richMenuID); !ok {
		writeLINEError(w, http.StatusBadRequest, "richmenu not found")
		return
	}
	if _, err := r.store.Update("line", "rich_menu_alias", resource.ID, map[string]any{"richMenuId": richMenuID}); err != nil {
		writeLINEError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeEmptyJSON(w, http.StatusOK)
}

func (r *routes) writeDeleteRichMenuAlias(w http.ResponseWriter, aliasID string) {
	resource, ok := r.findAlias(aliasID)
	if !ok {
		writeLINEError(w, http.StatusNotFound, "richmenu alias not found")
		return
	}
	r.store.Delete("line", "rich_menu_alias", resource.ID)
	writeEmptyJSON(w, http.StatusOK)
}

func (r *routes) writeGetRichMenuAlias(w http.ResponseWriter, aliasID string) {
	resource, ok := r.findAlias(aliasID)
	if !ok {
		writeLINEError(w, http.StatusNotFound, "richmenu alias not found")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"richMenuAliasId": resource.Data["richMenuAliasId"], "richMenuId": resource.Data["richMenuId"]})
}

func (r *routes) writeRichMenuAliasList(w http.ResponseWriter) {
	aliases := []map[string]any{}
	for _, resource := range r.store.List("line", "rich_menu_alias") {
		aliases = append(aliases, map[string]any{"richMenuAliasId": resource.Data["richMenuAliasId"], "richMenuId": resource.Data["richMenuId"]})
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"aliases": aliases})
}

func (r *routes) findAlias(aliasID string) (state.Resource, bool) {
	for _, resource := range r.store.List("line", "rich_menu_alias") {
		if resource.Data["richMenuAliasId"] == aliasID {
			return resource, true
		}
	}
	return state.Resource{}, false
}

func richMenuResponse(resource state.Resource) richMenuData {
	body := make(richMenuData, len(resource.Data)+1)
	for key, value := range resource.Data {
		body[key] = value
	}
	body["richMenuId"] = resource.ID
	delete(body, "hasImage")
	return body
}
