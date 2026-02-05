import { createFileRoute, Link, Outlet, redirect, useNavigate } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { useAuthStore } from "../../store/useAuthStore";

export const Route = createFileRoute("/examiner")({
    component: ExaminerLayout,
    beforeLoad: () => {
        const { token } = useAuthStore.getState();
        if (!token) {
            throw redirect({ to: "/login" });
        }
    },
});

function ExaminerLayout() {
    const { role, logout, isAuthenticated } = useAuthStore();
    const navigate = useNavigate();
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

    useEffect(() => {
        if (!isAuthenticated) {
            navigate({ to: "/login" as "/login" });
        }
    }, [isAuthenticated, navigate]);

    useEffect(() => {
        if (role && role !== "examiner" && role !== "admin") {
            // Optional: Strict redirect
        }
    }, [role]);

    // Close menu when route changes
    useEffect(() => {
        setIsMobileMenuOpen(false);
    }, [navigate]);

    return (
        <div className="min-h-screen bg-slate-50 font-sans text-slate-900">
            <header className="bg-white/80 backdrop-blur-md border-b border-slate-200/60 sticky top-0 z-50 px-4 md:px-6 py-3 flex items-center justify-between shadow-sm transition-all duration-300">
                {/* Left: Logo & Title */}
                <div className="flex items-center gap-3 md:w-1/4">
                    <div className="w-9 h-9 rounded-xl bg-gradient-to-tr from-rose-500 to-purple-600 flex items-center justify-center text-white font-bold text-xs shadow-lg shadow-purple-500/30 transform hover:scale-110 transition-transform duration-300 cursor-default">
                        BGK
                    </div>
                    <Link
                        to="/examiner"
                        search={{ tab: 'scoring' }}
                        activeOptions={{ exact: true }}
                        className="font-bold text-lg bg-clip-text text-transparent bg-gradient-to-r from-slate-800 to-slate-600 tracking-tight hover:from-purple-600 hover:to-rose-500 transition-all duration-300"
                    >
                        Cổng Giám Khảo
                    </Link>
                </div>

                {/* Center: Desktop Navigation Menu */}
                <nav className="hidden md:flex flex-1 justify-center">
                    <div className="flex items-center gap-4">
                        <Link
                            to="/examiner"
                            search={{ tab: 'scoring' }}
                            activeOptions={{ includeSearch: true }}
                            className="relative px-6 py-2.5 text-sm font-bold rounded-full transition-all duration-300 border border-transparent text-slate-500 hover:text-purple-600 hover:bg-purple-50 [&.active]:bg-white [&.active]:text-purple-600 [&.active]:border-purple-200 [&.active]:shadow-md shadow-sm"
                        >
                            Danh sách chấm thi
                        </Link>
                        <Link
                            to="/examiner"
                            search={{ tab: 'contestants' }}
                            className="relative px-6 py-2.5 text-sm font-bold rounded-full transition-all duration-300 border border-transparent text-slate-500 hover:text-purple-600 hover:bg-purple-50 [&.active]:bg-white [&.active]:text-purple-600 [&.active]:border-purple-200 [&.active]:shadow-md shadow-sm"
                        >
                            Danh sách thí sinh
                        </Link>
                        <Link
                            to="/examiner"
                            search={{ tab: 'approval' }}
                            className="relative px-6 py-2.5 text-sm font-bold rounded-full transition-all duration-300 border border-transparent text-slate-500 hover:text-purple-600 hover:bg-purple-50 [&.active]:bg-white [&.active]:text-purple-600 [&.active]:border-purple-200 [&.active]:shadow-md shadow-sm"
                        >
                            Duyệt Thí Sinh
                        </Link>
                    </div>
                </nav>

                {/* Right: Desktop User & Actions */}
                <div className="hidden md:flex items-center justify-end gap-5 w-1/4">
                    <div className="flex flex-col items-end">
                        <span className="text-xs font-bold text-slate-700 bg-slate-100 px-2 py-0.5 rounded-full border border-slate-200">
                            {role === 'examiner' ? 'Giám Khảo' : role}
                        </span>
                    </div>
                    <button
                        onClick={() => {
                            logout();
                            navigate({ to: "/login" as "/login" });
                        }}
                        className="group flex items-center gap-2 text-sm font-medium text-slate-500 hover:text-rose-600 transition-colors bg-white hover:bg-rose-50 px-4 py-2 rounded-xl border border-dashed border-slate-300 hover:border-rose-200"
                    >
                        <span>Đăng xuất</span>
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" className="w-4 h-4 group-hover:translate-x-1 transition-transform">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15m3 0l3-3m0 0l-3-3m3 3H9" />
                        </svg>
                    </button>
                </div>

                {/* Mobile Hamburger Button */}
                <button
                    className="md:hidden p-2 text-slate-600 hover:bg-slate-100 rounded-lg transition-colors"
                    onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                >
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" className="w-6 h-6">
                        {isMobileMenuOpen ? (
                            <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                        ) : (
                            <path strokeLinecap="round" strokeLinejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
                        )}
                    </svg>
                </button>

                {/* Mobile Menu Dropdown */}
                {isMobileMenuOpen && (
                    <div className="absolute top-full left-0 right-0 bg-white border-b border-slate-200 shadow-xl p-4 md:hidden flex flex-col gap-3 animate-in fade-in slide-in-from-top-4 duration-200">
                        <div className="flex items-center justify-between mb-2 pb-2 border-b border-slate-100">
                            <span className="text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Xin chào, {role === 'examiner' ? 'Giám Khảo' : role}
                            </span>
                        </div>
                        <Link
                            to="/examiner"
                            search={{ tab: 'scoring' }}
                            activeOptions={{ includeSearch: true }}
                            className="block px-4 py-3 rounded-xl bg-slate-50 text-slate-700 font-bold hover:bg-purple-50 hover:text-purple-600 transition-colors [&.active]:bg-purple-600 [&.active]:text-white shadow-sm"
                            onClick={() => setIsMobileMenuOpen(false)}
                        >
                            Danh sách chấm thi
                        </Link>
                        <Link
                            to="/examiner"
                            search={{ tab: 'contestants' }}
                            className="block px-4 py-3 rounded-xl bg-slate-50 text-slate-700 font-bold hover:bg-purple-50 hover:text-purple-600 transition-colors [&.active]:bg-purple-600 [&.active]:text-white shadow-sm"
                            onClick={() => setIsMobileMenuOpen(false)}
                        >
                            Danh sách thí sinh
                        </Link>
                        <Link
                            to="/examiner"
                            search={{ tab: 'approval' }}
                            className="block px-4 py-3 rounded-xl bg-slate-50 text-slate-700 font-bold hover:bg-purple-50 hover:text-purple-600 transition-colors [&.active]:bg-purple-600 [&.active]:text-white shadow-sm"
                            onClick={() => setIsMobileMenuOpen(false)}
                        >
                            Duyệt Thí Sinh
                        </Link>

                        <div className="my-2 border-t border-slate-100 pt-3">
                            <button
                                onClick={() => {
                                    logout();
                                    navigate({ to: "/login" as "/login" });
                                }}
                                className="w-full flex items-center justify-center gap-2 text-sm font-bold text-rose-600 bg-rose-50 hover:bg-rose-100 px-4 py-3 rounded-xl transition-colors"
                            >
                                <span>Đăng xuất</span>
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" className="w-4 h-4">
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15m3 0l3-3m0 0l-3-3m3 3H9" />
                                </svg>
                            </button>
                        </div>
                    </div>
                )}
            </header>

            <main className="max-w-7xl mx-auto p-4 md:p-6">
                <Outlet />
            </main>
        </div>
    );
}
