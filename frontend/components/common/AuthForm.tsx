"use client";

import { useRouter } from "next/navigation";
import { useDispatch } from "react-redux";
import { login, signup } from "../store/slice/auth";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Button } from "../ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel } from "../ui/form";

interface AuthFormProps {
  isLogin: boolean;
}

const formSchema = z.object({
  username: z.string(),
  email: z.string().email(),
  password: z.string().min(6),
});

const AuthForm: React.FC<AuthFormProps> = ({ isLogin }) => {
  const dispatch: any = useDispatch();
  const router = useRouter();
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
      email: "",
      password: "",
    },
  });

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    const { username, email, password } = values;
    if (isLogin) {
      await dispatch(login({ email, password }));
    } else {
      await dispatch(signup({ username, email, password }));
    }
    router.push("/");
  };

  return (
    <div className="flex min-h-full flex-1 flex-col justify-center px-6 py-12 lg:px-10">
      <div className="sm:mx-auto sm:w-full sm:max-w-sm">
        <img
          className="mx-auto h-10 w-auto"
          src="https://tailwindui.com/img/logos/mark.svg?color=indigo&shade=600"
          alt="Your Company"
        />
        <h2 className="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">
          {isLogin ? "Sign in to your account" : "Create an account"}
        </h2>
      </div>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="mt-10 grid grid-cols-1 gap-x-6 gap-y-4">
            {!isLogin && (
              <FormField
                name="username"
                control={form.control}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="block text-sm font-medium leading-6 text-gray-900">
                      Username
                    </FormLabel>
                    <FormControl>
                      <div className="mt-2">
                        <div className="flex rounded-md shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-600 sm:max-w-md">
                          <span className="flex select-none items-center pl-3 text-gray-500 sm:text-sm">
                            @
                          </span>
                          <input
                            {...field}
                            type="text"
                            id="username"
                            autoComplete="username"
                            className="block flex-1 border-0 bg-transparent py-1.5 pl-1 text-gray-900 outline-none placeholder:text-gray-400 focus:ring-0 sm:text-sm sm:leading-6"
                            placeholder="janesmith"
                          />
                        </div>
                      </div>
                    </FormControl>
                  </FormItem>
                )}
              />
            )}
            <FormField
              name="email"
              control={form.control}
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="block text-sm font-medium leading-6 text-gray-900">
                    Email Address
                  </FormLabel>
                  <FormControl>
                    <div className="mt-2">
                      <input
                        {...field}
                        id="email"
                        name="email"
                        type="email"
                        autoComplete="email"
                        className="px-2 block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                      />
                    </div>
                  </FormControl>
                </FormItem>
              )}
            />
            <FormField
              name="password"
              control={form.control}
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="block text-sm font-medium leading-6 text-gray-900">
                    Password
                  </FormLabel>
                  <FormControl>
                    <div className="mt-2">
                      <input
                        {...field}
                        id="password"
                        name="password"
                        type="password"
                        autoComplete="password"
                        className="px-2 block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                      />
                    </div>
                  </FormControl>
                </FormItem>
              )}
            />
          </div>
          <div className="mt-12 flex items-center justify-end gap-x-6">
            <Button
              type="submit"
              className="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
            >
              {isLogin ? "Login" : "Register"}
            </Button>
          </div>
        </form>
      </Form>
      <p className="mt-4 text-center text-sm text-gray-500">
        Not a member?{" "}
        <a
          href="/register"
          className="font-semibold leading-6 text-indigo-600 hover:text-indigo-500"
        >
          Register here
        </a>
      </p>
    </div>
  );
};

export default AuthForm;
