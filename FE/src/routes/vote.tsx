import { createFileRoute, Link, Outlet, useNavigate } from '@tanstack/react-router';
import { useState } from 'react';
import { useAuthStore } from '../store/useAuthStore';

function VoteLayout() {
    const navigate = useNavigate();
    const { role, isAuthenticated, logout } = useAuthStore();

    const [isMenuOpen, setIsMenuOpen] = useState(false);

    return (
        <div className="min-h-screen relative font-sans text-gray-900">
            {/* Fixed Background */}
            <div className="fixed inset-0 z-0" style={{
                backgroundImage: `url('/background_vote.jpeg')`,
                backgroundSize: 'cover',
                backgroundPosition: 'center',
            }}>
                <div className="absolute inset-0 bg-black/30"></div>
            </div>

            {/* Content Wrapper */}
            <div className="relative z-10 flex flex-col min-h-screen">
                {/* Marquee Banner */}
                <div className="bg-white text-gray-900 py-2 overflow-hidden border-b border-gray-200">
                    <div className="animate-marquee whitespace-nowrap">
                        <span className="mx-4 font-semibold text-sm md:text-base">🎉 Người đẹp được yêu thích nhất - vào tháng TOP 10 chung cuộc</span>
                        <span className="mx-4 font-semibold text-sm md:text-base">🎉 Người đẹp được yêu thích nhất - vào tháng TOP 10 chung cuộc</span>
                        <span className="mx-4 font-semibold text-sm md:text-base">🎉 Người đẹp được yêu thích nhất - vào tháng TOP 10 chung cuộc</span>
                        <span className="mx-4 font-semibold text-sm md:text-base">🎉 Người đẹp được yêu thích nhất - vào tháng TOP 10 chung cuộc</span>
                    </div>
                </div>

                {/* Top Navigation Bar */}
                <nav className="bg-black border-b border-gray-800 relative z-50">
                    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                        <div className="flex justify-between items-center h-14">
                            {/* Logo */}
                            <div className="flex items-center space-x-2">
                                <img src="/logoM.png" alt="Logo" className="h-8 w-auto" />
                                <h1 className="text-xl font-bold text-white tracking-wider">EVENTISTA</h1>
                            </div>

                            {/* Mobile Menu Button */}
                            <button
                                onClick={() => setIsMenuOpen(!isMenuOpen)}
                                className="md:hidden text-white p-2 rounded-md hover:bg-white/10 transition"
                            >
                                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    {isMenuOpen ? (
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                    ) : (
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                                    )}
                                </svg>
                            </button>

                            {/* Desktop Right Menu */}
                            <div className="hidden md:flex items-center space-x-6">
                                <a href="#" className="text-white/90 hover:text-white text-sm transition">Giới thiệu Eventista</a>
                                <span className="text-white/40">|</span>
                                <a href="#" className="text-white/90 hover:text-white text-sm transition">Tất cả sự kiện</a>
                                <span className="text-white/40">|</span>
                                {isAuthenticated ? (
                                    <div className="flex items-center space-x-3">
                                        <span className="text-white/90 text-sm">Xin chào, <span className="font-semibold">{role}</span></span>
                                        <button
                                            onClick={logout}
                                            className="text-white/90 hover:text-white text-sm transition"
                                        >
                                            Đăng xuất
                                        </button>
                                    </div>
                                ) : (
                                    <button
                                        onClick={() => navigate({ to: '/vote/auth' as '/vote/auth' })}
                                        className="flex items-center space-x-1 text-white/90 hover:text-white text-sm transition"
                                    >
                                        <span>Đăng nhập</span>
                                        <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                                            <path fillRule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clipRule="evenodd" />
                                        </svg>
                                    </button>
                                )}
                            </div>
                        </div>
                    </div>

                    {/* Mobile Menu Overlay */}
                    {isMenuOpen && (
                        <div className="md:hidden bg-black/95 border-b border-gray-800 absolute w-full left-0 top-14 shadow-xl">
                            <div className="px-4 pt-2 pb-4 space-y-1">
                                <div className="border-b border-white/10 pb-2 mb-2">
                                    <Link to="/vote" className="block px-3 py-2 text-base font-medium text-white hover:bg-white/10 rounded-md" onClick={() => setIsMenuOpen(false)}>Trang bình chọn</Link>
                                    <Link to="/vote/bang-xep-hang" className="block px-3 py-2 text-base font-medium text-white hover:bg-white/10 rounded-md" onClick={() => setIsMenuOpen(false)}>Bảng xếp hạng</Link>
                                    <Link to="/vote/danh-sach-thi-sinh" className="block px-3 py-2 text-base font-medium text-white hover:bg-white/10 rounded-md" onClick={() => setIsMenuOpen(false)}>Danh sách thí sinh</Link>
                                </div>
                                <a href="#" className="block px-3 py-2 text-sm text-white/70 hover:text-white hover:bg-white/5 rounded-md">Giới thiệu Eventista</a>
                                <a href="#" className="block px-3 py-2 text-sm text-white/70 hover:text-white hover:bg-white/5 rounded-md">Tất cả sự kiện</a>
                                {isAuthenticated ? (
                                    <div className="mt-4 border-t border-white/10 pt-4 px-3">
                                        <div className="text-white/90 text-sm mb-2">Xin chào, <span className="font-semibold">{role}</span></div>
                                        <button
                                            onClick={() => { logout(); setIsMenuOpen(false); }}
                                            className="w-full text-left text-red-400 hover:text-red-300 text-sm transition py-2"
                                        >
                                            Đăng xuất
                                        </button>
                                    </div>
                                ) : (
                                    <button
                                        onClick={() => { navigate({ to: '/vote/auth' as '/vote/auth' }); setIsMenuOpen(false); }}
                                        className="mt-4 w-full flex items-center justify-center space-x-2 bg-white text-black px-4 py-2 rounded-md font-medium hover:bg-gray-100 transition"
                                    >
                                        <span>Đăng nhập</span>
                                        <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                                            <path fillRule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clipRule="evenodd" />
                                        </svg>
                                    </button>
                                )}
                            </div>
                        </div>
                    )}
                </nav>

                {/* Sub Navigation Bar - Hidden on Mobile */}
                <nav className="hidden md:block bg-black/80 border-b border-gray-800">
                    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                        <div className="flex items-center justify-center space-x-8 h-12">
                            <Link to="/vote" className="text-white/80 hover:text-white text-sm transition" activeProps={{ className: '!font-semibold !text-white !border-b-2 !border-white !pb-3' }}>Trang bình chọn</Link>
                            <span className="text-white/20">|</span>
                            <Link to="/vote/bang-xep-hang" className="text-white/80 hover:text-white text-sm transition" activeProps={{ className: '!font-semibold !text-white !border-b-2 !border-white !pb-3' }}>Bảng xếp hạng</Link>
                            <span className="text-white/20">|</span>
                            <Link to="/vote/danh-sach-thi-sinh" className="text-white/80 hover:text-white text-sm transition" activeProps={{ className: '!font-semibold !text-white !border-b-2 !border-white !pb-3' }}>Danh sách thí sinh</Link>
                        </div>
                    </div>
                </nav>

                {/* Main Content Area */}
                <main className="flex-1 w-full">
                    <Outlet />
                </main>

                {/* Footer */}
                <footer className="bg-black/80 backdrop-blur-md border-t border-white/10 py-6 mt-auto">
                    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center text-white/60 text-sm">
                        © 2026 Cuộc Thi Hoa Hậu Việt Nam - EVENTISTA
                    </div>
                </footer>
            </div>
            {/* End Content Wrapper */}
        </div>
    );
}

export const Route = createFileRoute('/vote')({
    component: VoteLayout,
});
