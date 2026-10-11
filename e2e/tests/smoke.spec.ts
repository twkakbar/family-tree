import { expect, test } from "@playwright/test";

test.describe("public smoke", () => {
  test("explore home shows royal brand", async ({ page }) => {
    await page.goto("/");
    await expect(page.getByText("Kesultanan Aceh").first()).toBeVisible();
    await expect(page.getByRole("link", { name: "About" })).toBeVisible();
  });

  test("about page loads", async ({ page }) => {
    await page.goto("/about");
    await expect(page.getByText(/Aceh|Kesultanan|family/i).first()).toBeVisible();
  });

  test("login page shows sign-in form", async ({ page }) => {
    await page.goto("/login");
    await expect(page.getByLabel("Email")).toBeVisible();
    await expect(page.getByLabel("Password")).toBeVisible();
    await expect(page.getByRole("button", { name: /sign in/i })).toBeVisible();
  });
});

test.describe("api probes", () => {
  test("healthz is ok", async ({ request }) => {
    const res = await request.get("/healthz");
    expect(res.ok()).toBeTruthy();
    const body = await res.json();
    expect(body.status).toBe("ok");
  });

  test("readyz is ok", async ({ request }) => {
    const res = await request.get("/readyz");
    expect(res.ok()).toBeTruthy();
    const body = await res.json();
    expect(body.status).toBe("ok");
  });

  test("public meta returns royal family", async ({ request }) => {
    const res = await request.get("/api/v1/public/meta");
    const body = await res.text();
    expect(res.ok(), body).toBeTruthy();
    const envelope = JSON.parse(body) as { success?: boolean; data?: { name?: string } };
    expect(envelope.success).toBe(true);
    expect(envelope.data?.name).toBeTruthy();
  });

  test("security headers include CSP", async ({ request }) => {
    const res = await request.get("/");
    expect(res.ok()).toBeTruthy();
    const csp = res.headers()["content-security-policy"];
    expect(csp).toBeTruthy();
    expect(csp).toContain("frame-ancestors 'none'");
    expect(res.headers()["x-frame-options"]?.toLowerCase()).toBe("deny");
  });
});

test.describe("auth smoke", () => {
  test("register → dashboard shell", async ({ page }) => {
    const stamp = Date.now();
    const email = `e2e-${stamp}@example.com`;
    const password = "password123";

    await page.goto("/register");
    await page.getByLabel("Full name").fill("E2E User");
    await page.getByLabel("Email").fill(email);
    await page.getByLabel("Password", { exact: true }).fill(password);
    await page.getByLabel("Confirm password").fill(password);
    await page.getByRole("button", { name: /create account/i }).click();

    await expect(page).toHaveURL(/\/(dashboard|request-access|studio)/, { timeout: 30_000 });
  });
});
