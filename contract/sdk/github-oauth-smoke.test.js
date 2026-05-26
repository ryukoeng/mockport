"use strict";

async function runGitHubOAuthSmoke(options) {
  const baseURL = new URL("/github", options.baseURL).toString().replace(/\/$/, "");
  const redirectURI = "http://app.local/callback";
  const authURL = new URL(`${baseURL}/login/oauth/authorize`);
  authURL.searchParams.set("client_id", "mockport_github_client");
  authURL.searchParams.set("redirect_uri", redirectURI);
  authURL.searchParams.set("state", "state-123");
  authURL.searchParams.set("scope", "read:user user:email read:org");

  const auth = await fetch(authURL, { redirect: "manual" });
  assertEqual(auth.status, 302, "authorize status");
  const location = auth.headers.get("location");
  if (!location) {
    throw new Error("authorize redirect missing Location header");
  }
  const redirect = new URL(location);
  assertEqual(`${redirect.origin}${redirect.pathname}`, redirectURI, "redirect uri");
  assertEqual(redirect.searchParams.get("state"), "state-123", "redirect state");
  const code = redirect.searchParams.get("code");
  if (!code) {
    throw new Error("authorize redirect missing code");
  }

  const token = await fetch(`${baseURL}/login/oauth/access_token`, {
    method: "POST",
    headers: {
      "accept": "application/json",
      "content-type": "application/x-www-form-urlencoded",
    },
    body: new URLSearchParams({
      client_id: "mockport_github_client",
      client_secret: "mockport_github_secret",
      code,
      redirect_uri: redirectURI,
    }),
  });
  assertEqual(token.status, 200, "token status");
  const tokenBody = await token.json();
  assertEqual(tokenBody.token_type, "bearer", "token type");
  assertEqual(tokenBody.scope, "read:user user:email read:org", "token scope");
  if (!tokenBody.access_token) {
    throw new Error("token response missing access_token");
  }

  const user = await githubJSON(`${baseURL}/user`, tokenBody.access_token);
  assertEqual(user.login, "mockport-user", "user login");
  assertEqual(user.email, "mockport@example.test", "user email");

  const emails = await githubJSON(`${baseURL}/user/emails`, tokenBody.access_token);
  assertEqual(emails[0].email, "mockport@example.test", "primary email");
  assertEqual(emails[0].primary, true, "primary email flag");

  const orgs = await githubJSON(`${baseURL}/user/orgs`, tokenBody.access_token);
  assertEqual(orgs[0].login, "mockport-org", "org login");

  return {
    provider: "github-oauth",
    baseURL: options.baseURL,
    status: "client-ok",
    code,
    token: tokenBody.access_token,
    user: user.login,
    email: emails[0].email,
    org: orgs[0].login,
  };
}

async function githubJSON(url, token) {
  const response = await fetch(url, {
    headers: {
      "accept": "application/vnd.github+json",
      "authorization": `Bearer ${token}`,
    },
  });
  if (!response.ok) {
    throw new Error(`${url} returned ${response.status}: ${await response.text()}`);
  }
  return response.json();
}

function assertEqual(got, want, label) {
  if (got !== want) {
    throw new Error(`${label}: got ${got}, want ${want}`);
  }
}

module.exports = { runGitHubOAuthSmoke };
