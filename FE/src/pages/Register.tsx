import { Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useForm } from "react-hook-form";

import { apiRequest } from "../lib/api";

const REGISTER_MUTATION = `
mutation Register($email: String!, $password: String!) {
  register(email: $email, password: $password) {
    token
  }
}
`;

interface RegisterFormData {
    email: string;
    password: string;
    confirmPassword: string;
}

export const Register = () => {
    const { register, handleSubmit } = useForm<RegisterFormData>();
    const [error, setError] = useState("");
    const navigate = useNavigate();

    const onSubmit = async (data: RegisterFormData) => {
        setError("");
        if (data.password !== data.confirmPassword) {
            setError("Passwords do not match");
            return;
        }

        try {
            const result = await apiRequest(REGISTER_MUTATION, { email: data.email, password: data.password });
            const token = result.register.token;
            localStorage.setItem("token", token);

            // Redirect to dashboard
            navigate({ to: "/contestant" });
            window.location.href = "/contestant";
        } catch (err: any) {
            setError(err.message || "Failed to register");
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8 bg-cover bg-center" style={{ backgroundImage: "url('/login.jpeg')" }}>
            <div className="max-w-4xl w-full bg-white/90 backdrop-blur-sm rounded-2xl shadow-2xl overflow-hidden flex shadow-slate-200/50">
                {/* Left Side: Image */}
                <div className="hidden md:block w-1/2 bg-cover bg-center" style={{ backgroundImage: "url('/1.jpeg')" }}>
                </div>

                {/* Right Side: Form */}
                <div className="w-full md:w-1/2 p-8 md:p-12 lg:p-16 flex flex-col justify-center relative">
                    <div className="sm:mx-auto sm:w-full sm:max-w-md">
                        <h2 className="mt-2 text-center text-3xl font-extrabold tracking-tight text-gray-900">
                            Tham Gia Cuộc Thi
                        </h2>
                        <p className="mt-2 text-center text-sm text-gray-600">
                            Tạo tài khoản để tham gia cuộc thi
                        </p>
                    </div>

                    <form className="mt-8 space-y-6" onSubmit={handleSubmit(onSubmit)}>
                        <div className="space-y-4">
                            <div>
                                <label htmlFor="email" className="block text-sm font-medium leading-6 text-gray-900">
                                    Email
                                </label>
                                <div className="mt-1">
                                    <input
                                        id="email"
                                        type="email"
                                        required
                                        className="block w-full rounded-lg border-0 py-2.5 text-gray-900 chat-shadow ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-pink-500 sm:text-sm sm:leading-6 px-4 transition-all duration-200"
                                        placeholder="you@example.com"
                                        {...register("email", { required: true })}
                                    />
                                </div>
                            </div>

                            <div>
                                <label htmlFor="password" className="block text-sm font-medium leading-6 text-gray-900">
                                    Password
                                </label>
                                <div className="mt-1">
                                    <input
                                        id="password"
                                        type="password"
                                        required
                                        className="block w-full rounded-lg border-0 py-2.5 text-gray-900 chat-shadow ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-pink-500 sm:text-sm sm:leading-6 px-4 transition-all duration-200"
                                        placeholder="••••••••"
                                        {...register("password", { required: true })}
                                    />
                                </div>
                            </div>

                            <div>
                                <label htmlFor="confirmPassword" className="block text-sm font-medium leading-6 text-gray-900">
                                    Confirm Password
                                </label>
                                <div className="mt-1">
                                    <input
                                        id="confirmPassword"
                                        type="password"
                                        required
                                        className="block w-full rounded-lg border-0 py-2.5 text-gray-900 chat-shadow ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-pink-500 sm:text-sm sm:leading-6 px-4 transition-all duration-200"
                                        placeholder="••••••••"
                                        {...register("confirmPassword", { required: true })}
                                    />
                                </div>
                            </div>
                        </div>

                        {error && (
                            <div className="rounded-md bg-red-50 p-4">
                                <div className="flex">
                                    <div className="ml-3">
                                        <h3 className="text-sm font-medium text-red-800">Registration Failed</h3>
                                        <div className="mt-2 text-sm text-red-700">
                                            <p>{error}</p>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        )}

                        <div>
                            <button
                                type="submit"
                                className="flex w-full justify-center rounded-lg
bg-gradient-to-r from-rose-400 via-pink-400 to-purple-400
px-3 py-3 text-sm font-semibold leading-6 text-white
shadow-[0_15px_35px_rgba(200,150,200,0.45)]
hover:from-rose-300 hover:via-pink-300 hover:to-purple-300
focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-400
transition-all duration-300 transform hover:-translate-y-0.5

                                        "
                            >
                                Đăng ký
                            </button>
                        </div>
                    </form>

                    <p className="mt-10 text-center text-sm text-gray-500">
                        Đã có tài khoản? {" "}
                        <Link to="/login" className="font-semibold leading-6 text-purple-600 hover:text-purple-500 transition-colors">
                            Đăng nhập
                        </Link>
                    </p>
                </div>
            </div>
        </div>
    );
};
