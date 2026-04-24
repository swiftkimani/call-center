import { LoginForm } from "@/app/login/login-form";
import { BrandLogo } from "@/components/layout/brand-logo";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export function LoginScreen() {
  return (
    <main className="mx-auto flex min-h-screen max-w-6xl items-center px-4 py-10 md:px-6">
      <div className="grid w-full gap-10 lg:grid-cols-[0.95fr_0.7fr]">
        <section className="space-y-6">
          <BrandLogo />
          <div className="max-w-xl space-y-4">
            <h1 className="text-4xl font-semibold tracking-tight text-foreground">Sign in and get straight to the call floor.</h1>
            <p className="text-base leading-7 text-muted-foreground">
              The login stays intentionally light. Authenticate, enter the workspace, and leave the noise outside the product.
            </p>
          </div>
        </section>

        <Card className="mx-auto w-full max-w-md">
          <CardHeader>
            <CardTitle className="text-2xl">Sign in</CardTitle>
            <CardDescription>Use your assigned call center account.</CardDescription>
          </CardHeader>
          <CardContent>
            <LoginForm />
          </CardContent>
        </Card>
      </div>
    </main>
  );
}
