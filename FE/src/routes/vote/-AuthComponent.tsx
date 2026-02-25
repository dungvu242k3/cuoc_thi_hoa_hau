import { useNavigate } from "@tanstack/react-router"; // or your router
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useAuthStore } from "../../store/useAuthStore";

// Simple Tab Component
function Tabs({ activeTab, onTabChange }: { activeTab: string, onTabChange: (tab: string) => void }) {
    return (
        <div className="flex border-b border-gray-200 mb-6">
            <button
                className={`flex-1 py-3 text-sm font-medium ${activeTab === 'login' ? 'text-pink-600 border-b-2 border-pink-600' : 'text-gray-500 hover:text-gray-700'}`}
                onClick={() => onTabChange('login')}
            >
                Đăng Nhập
            </button>
            <button
                className={`flex-1 py-3 text-sm font-medium ${activeTab === 'register' ? 'text-pink-600 border-b-2 border-pink-600' : 'text-gray-500 hover:text-gray-700'}`}
                onClick={() => onTabChange('register')}
            >
                Đăng Ký
            </button>
        </div>
    );
}

// Rename to avoid conflict if needed, or keeping it is fine but make sure route.ts imports it correctly
export function VoteAuthComponent() {
    const [activeTab, setActiveTab] = useState<'login' | 'register'>('login');
    const navigate = useNavigate();
    // const setAuth = useAuthStore((state) => state.setAuth);

    const [error, setError] = useState<string | null>(null);

    const { register: loginRegister, handleSubmit: handleLoginSubmit, formState: { errors: loginErrors, isSubmitting: isLoginSubmitting } } = useForm<any>();
    const { register: registerRegister, handleSubmit: handleRegisterSubmit, formState: { errors: registerErrors, isSubmitting: isRegisterSubmitting } } = useForm<any>();

    const onLogin = async (data: any) => {
        setError(null);
        try {
            const res = await fetch(`http://${window.location.hostname}:8080/api/login-audience`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data),
            });

            if (!res.ok) {
                throw new Error(await res.text() || 'Đăng nhập thất bại');
            }

            const body: any = await res.json(); // { token, userId, role }

            // Update Store
            useAuthStore.getState().login(body.token);

            // Redirect
            navigate({ to: '/vote' as '/vote' });
        } catch (err: any) {
            setError(err.message || 'Đăng nhập thất bại');
        }
    };

    const onRegister = async (data: any) => {
        setError(null);
        try {
            const res = await fetch(`http://${window.location.hostname}:8080/api/register-audience`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data),
            });

            if (!res.ok) {
                throw new Error(await res.text() || 'Đăng ký thất bại');
            }

            const body: any = await res.json();
            // Update Store
            useAuthStore.getState().login(body.token);

            // Redirect
            navigate({ to: '/vote' });
        } catch (err: any) {
            setError(err.message || 'Đăng ký thất bại');
        }
    };

    return (
        <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
            <div className="sm:mx-auto sm:w-full sm:max-w-md">
                <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
                    Cổng Bình Chọn Khán Giả
                </h2>
                <p className="mt-2 text-center text-sm text-gray-600">
                    Đăng nhập để bình chọn cho thí sinh yêu thích
                </p>
            </div>

            <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
                <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
                    <Tabs activeTab={activeTab} onTabChange={(t) => setActiveTab(t as any)} />

                    {error && (
                        <div className="mb-4 bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded relative" role="alert">
                            <span className="block sm:inline">{error}</span>
                        </div>
                    )}

                    {activeTab === 'login' && (
                        <form className="space-y-6" onSubmit={handleLoginSubmit(onLogin)}>
                            <div>
                                <label htmlFor="email" className="block text-sm font-medium text-gray-700">Email / Số điện thoại</label>
                                <div className="mt-1">
                                    <input
                                        id="email"
                                        type="text"
                                        {...loginRegister('email', { required: 'Vui lòng nhập email' })}
                                        className="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-pink-500 focus:border-pink-500 sm:text-sm"
                                    />
                                    {loginErrors['email'] && <p className="mt-1 text-sm text-red-600">{loginErrors['email'].message as string}</p>}
                                </div>
                            </div>

                            <div>
                                <label htmlFor="password" className="block text-sm font-medium text-gray-700">Mật khẩu</label>
                                <div className="mt-1">
                                    <input
                                        id="password"
                                        type="password"
                                        {...loginRegister('password', { required: 'Vui lòng nhập mật khẩu' })}
                                        className="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-pink-500 focus:border-pink-500 sm:text-sm"
                                    />
                                    {loginErrors['password'] && <p className="mt-1 text-sm text-red-600">{loginErrors['password'].message as string}</p>}
                                </div>
                            </div>

                            <div>
                                <button
                                    type="submit"
                                    disabled={isLoginSubmitting}
                                    className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-pink-600 hover:bg-pink-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-pink-500 disabled:opacity-50"
                                >
                                    {isLoginSubmitting ? 'Đang xử lý...' : 'Đăng Nhập'}
                                </button>
                            </div>
                        </form>
                    )}

                    {activeTab === 'register' && (
                        <form className="space-y-6" onSubmit={handleRegisterSubmit(onRegister)}>
                            <div>
                                <label className="block text-sm font-medium text-gray-700">Email</label>
                                <div className="mt-1">
                                    <input
                                        type="email"
                                        {...registerRegister('email', { required: 'Vui lòng nhập email', pattern: { value: /^\S+@\S+$/i, message: 'Email không hợp lệ' } })}
                                        className="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-pink-500 focus:border-pink-500 sm:text-sm"
                                    />
                                    {registerErrors['email'] && <p className="mt-1 text-sm text-red-600">{registerErrors['email'].message as string}</p>}
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700">Mật khẩu</label>
                                <div className="mt-1">
                                    <input
                                        type="password"
                                        {...registerRegister('password', {
                                            required: 'Vui lòng nhập mật khẩu',
                                            minLength: { value: 8, message: 'Mật khẩu tối thiểu 8 ký tự' }
                                        })}
                                        className="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-pink-500 focus:border-pink-500 sm:text-sm"
                                    />
                                    {registerErrors['password'] && <p className="mt-1 text-sm text-red-600">{registerErrors['password'].message as string}</p>}
                                    <p className="mt-1 text-xs text-gray-500">Yêu cầu: Chữ hoa, chữ thường, số và ký tự đặc biệt.</p>
                                </div>
                            </div>

                            <div>
                                <button
                                    type="submit"
                                    disabled={isRegisterSubmitting}
                                    className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-pink-600 hover:bg-pink-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-pink-500 disabled:opacity-50"
                                >
                                    {isRegisterSubmitting ? 'Đang xử lý...' : 'Đăng Ký'}
                                </button>
                            </div>
                        </form>
                    )}
                </div>
            </div>
        </div>
    );
}
