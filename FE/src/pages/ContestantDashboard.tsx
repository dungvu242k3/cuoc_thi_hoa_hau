import { useNavigate } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { apiRequest } from "../lib/api";

// --- Icons ---
const HomeIcon = () => (<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5"><path strokeLinecap="round" strokeLinejoin="round" d="M2.25 12l8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25" /></svg>);
const PhotoIcon = () => (<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5"><path strokeLinecap="round" strokeLinejoin="round" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 001.5-1.5V6a1.5 1.5 0 00-1.5-1.5H3.75A1.5 1.5 0 002.25 6v12a1.5 1.5 0 001.5 1.5z" /></svg>);
const ChartBarIcon = () => (<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5"><path strokeLinecap="round" strokeLinejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" /></svg>);
const UserIcon = () => (<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5"><path strokeLinecap="round" strokeLinejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" /></svg>);
const LogoutIcon = () => (<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5"><path strokeLinecap="round" strokeLinejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75" /></svg>);
const BellIcon = () => (<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-6 h-6"><path strokeLinecap="round" strokeLinejoin="round" d="M14.857 17.082a23.848 23.848 0 005.454-1.31A8.967 8.967 0 0118 9.75v-.7V9A6 6 0 006 9v.75a8.967 8.967 0 01-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 01-5.714 0m5.714 0a3 3 0 11-5.714 0" /></svg>);

const MY_PROFILE_QUERY = `
query MyProfile {
  myProfile {
        id
        sbd
        status
    personalInfo { fullName email dob nationality identityCard phone address job }
    physicalInfo { height weight measurements }
    skillEducation { educationLevel languages skills }
    portfolio { avatarUrl galleryUrls introduction socialLinks }
    }
}
`;

// ... existing imports

const MY_FEEDBACKS_QUERY = `
query MyFeedbacks {
    myFeedbacks {
        id
        title
        content
        type
        status
        reply
        createdAt
    }
}
`;

const SEND_FEEDBACK_MUTATION = `
mutation SendFeedback($input: CreateFeedbackInput!) {
    sendFeedback(input: $input)
}
`;

const MY_SCORES_QUERY = `
query MyScores {
    myScores {
        id
        roundId
        sbd
        totalScore
        criteriaScores
        comment
        createdAt
    }
}
`;

// ... existing queries

// Icons
const ChatIcon = () => (
    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-6 h-6">
        <path strokeLinecap="round" strokeLinejoin="round" d="M8.625 12a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H8.25m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H12.375m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0h-.375M21 12c0 4.556-4.03 8.25-9 8.25a9.764 9.764 0 01-2.555-.337A5.972 5.972 0 015.41 20.97a5.969 5.969 0 01-.474-.065 4.48 4.48 0 00.978-2.025c.09-.457-.133-.901-.467-1.226C3.93 16.178 3 14.189 3 12c0-4.556 4.03-8.25 9-8.25s9 3.694 9 8.25z" />
    </svg>
);

// ... existing icons

const UPDATE_PROFILE_MUTATION = `
mutation UpdateContestantProfile($input: UpdateContestantInput!) {
    updateContestantProfile(input: $input) {
        id
        sbd
        status
        personalInfo { fullName email dob nationality identityCard phone address job }
        physicalInfo { height weight measurements }
        skillEducation { educationLevel languages skills }
        portfolio { avatarUrl galleryUrls introduction socialLinks }
    }
}
`;

export const ContestantDashboard = () => {
    const navigate = useNavigate();
    const [profile, setProfile] = useState<any>(null);
    const [loading, setLoading] = useState(true);
    const [activeTab, setActiveTab] = useState("overview");
    const [uploading, setUploading] = useState(false);

    // Feedback State
    const [feedbacks, setFeedbacks] = useState<any[]>([]);
    const [isFeedbackModalOpen, setIsFeedbackModalOpen] = useState(false);
    const [feedbackTitle, setFeedbackTitle] = useState("");
    const [feedbackContent, setFeedbackContent] = useState("");
    const [feedbackType, setFeedbackType] = useState("support");
    const [loadingFeedbacks, setLoadingFeedbacks] = useState(false);

    useEffect(() => {
        const token = localStorage.getItem("token");
        if (!token) {
            navigate({ to: "/login" });
            return;
        }

        apiRequest(MY_PROFILE_QUERY, {}, token)
            .then((data: any) => {
                if (data && data.myProfile) {
                    setProfile(data.myProfile);
                } else {
                    navigate({ to: "/login" });
                }
            })
            .catch(() => navigate({ to: "/login" }))
            .finally(() => setLoading(false));
    }, [navigate]);

    // Fetch Feedbacks when tab changes
    useEffect(() => {
        if (activeTab === 'notifications') {
            fetchFeedbacks();
        }
    }, [activeTab]);

    const fetchFeedbacks = () => {
        const token = localStorage.getItem("token");
        setLoadingFeedbacks(true);
        apiRequest(MY_FEEDBACKS_QUERY, {}, token)
            .then((data: any) => {
                if (data && data.myFeedbacks) {
                    setFeedbacks(data.myFeedbacks);
                }
            })
            .catch(console.error)
            .finally(() => setLoadingFeedbacks(false));
    };

    const handleSendFeedback = async (e: React.FormEvent) => {
        e.preventDefault();
        const token = localStorage.getItem("token");

        if (!feedbackTitle || !feedbackContent) {
            alert("Vui lòng điền đầy đủ tiêu đề và nội dung.");
            return;
        }

        try {
            await apiRequest(SEND_FEEDBACK_MUTATION, {
                input: {
                    title: feedbackTitle,
                    content: feedbackContent,
                    type: feedbackType
                }
            }, token);
            alert("Gửi phản hồi thành công! Ban tổ chức sẽ sớm phản hồi bạn.");
            setIsFeedbackModalOpen(false);
            setFeedbackTitle("");
            setFeedbackContent("");
            setFeedbackType("support");
            if (activeTab === 'notifications') {
                fetchFeedbacks(); // Refresh list immediately
            } else {
                setActiveTab('notifications'); // Switch to tab to see it
            }
        } catch (err: any) {
            alert("Lỗi: " + err.message);
        }
    };

    // Score State
    const [scores, setScores] = useState<any[]>([]);
    const [loadingScores, setLoadingScores] = useState(false);

    useEffect(() => {
        if (activeTab === 'stats') {
            const token = localStorage.getItem("token");
            setLoadingScores(true);
            apiRequest(MY_SCORES_QUERY, {}, token)
                .then((data: any) => {
                    if (data && data.myScores) {
                        setScores(data.myScores);
                    }
                })
                .catch(console.error)
                .finally(() => setLoadingScores(false));
        }
    }, [activeTab]);

    const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>, fieldName: string) => {
        const file = e.target.files?.[0];
        if (!file) return;

        // Validation
        const validTypes = ['image/jpeg', 'image/png', 'image/webp'];
        if (!validTypes.includes(file.type)) {
            alert("Chỉ chấp nhận định dạng ảnh JPG, PNG, hoặc WebP.");
            return;
        }
        if (file.size > 5 * 1024 * 1024) {
            alert("Dung lượng ảnh không được vượt quá 5MB.");
            return;
        }

        setUploading(true);
        const formData = new FormData();
        formData.append("file", file);

        try {
            const token = localStorage.getItem("token");
            let apiUrl = import.meta.env.VITE_API_URL || "http://localhost:8080";
            if (apiUrl.endsWith("/query")) apiUrl = apiUrl.slice(0, -6);

            const response = await fetch(`${apiUrl}/upload`, {
                method: "POST",
                headers: { "Authorization": `Bearer ${token}` },
                body: formData
            });

            if (!response.ok) throw new Error("Upload failed");

            const data = await response.json() as { url: string };

            // Update local profile state immediately for preview
            setProfile((prev: any) => {
                if (!prev) return prev;
                const newPortfolio = { ...(prev.portfolio || {}), [fieldName]: data.url };
                return { ...prev, portfolio: newPortfolio };
            });

        } catch (error: any) {
            console.error("Upload error:", error);
            alert("Lỗi tải ảnh lên.");
        } finally {
            setUploading(false);
        }
    };

    const handleLogout = () => {
        localStorage.removeItem("token");
        navigate({ to: "/login" });
    }

    if (loading) return (
        <div className="flex min-h-screen items-center justify-center bg-white">
            <div className="flex flex-col items-center animate-pulse">
                <div className="h-12 w-12 bg-blue-900 rounded-full mb-4"></div>
                <div className="text-gray-400 font-light tracking-widest uppercase">Loading Portal...</div>
            </div>
        </div>
    );

    const NavItem = ({ tab, label, icon: Icon }: { tab: string, label: string, icon: any }) => (
        <button
            onClick={() => setActiveTab(tab)}
            className={`
                group relative flex flex-col items-center justify-center px-4 h-full transition-all duration-300
                ${activeTab === tab
                    ? 'text-yellow-400 scale-105'
                    : 'text-white hover:text-yellow-200'
                }
            `}
        >
            <div className={`mb-1 transition-all duration-300 ${activeTab === tab ? '-translate-y-1' : 'group-hover:-translate-y-1'}`}>
                <Icon className="w-6 h-6" />
            </div>
            <span className={`text-[10px] uppercase tracking-widest font-bold transition-all duration-300 ${activeTab === tab ? 'opacity-100' : 'opacity-90 group-hover:opacity-100'}`}>{label}</span>
            <span className={`absolute bottom-0 w-8 h-1 rounded-t-full bg-gradient-to-t from-yellow-400 to-transparent shadow-[0_-2px_6px_rgba(250,204,21,0.4)] transition-all duration-300 ${activeTab === tab ? 'scale-100 opacity-100' : 'scale-0 opacity-0'}`}></span>
        </button>
    );

    const getFullImageUrl = (path: string) => {
        if (!path) return "";
        if (path.startsWith("http") || path.startsWith("data:")) return path;
        const baseUrl = (import.meta.env.VITE_API_URL || "http://localhost:8080").replace(/\/query\/?$/, "");
        return `${baseUrl}${path}`;
    };

    return (
        <div className="min-h-screen font-sans text-slate-800 relative overflow-x-hidden">
            <div className="fixed inset-0 z-0 bg-blue-950/40 backdrop-blur-[2px]" style={{ backgroundImage: "url('/bg-dashboard.jpg')", backgroundSize: 'cover' }}></div>

            <div className="relative z-10 flex flex-col min-h-screen">
                {/* Header */}
                <header className="sticky top-0 z-50 bg-blue-950/50 backdrop-blur-xl border-b border-white/10 shadow-lg">
                    <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
                        <div className="flex justify-between items-center h-20">
                            <div className="flex items-center space-x-3 group">
                                <div className="h-11 w-11 bg-gradient-to-br from-yellow-400 to-yellow-600 text-blue-950 rounded-xl flex items-center justify-center font-serif font-bold cursor-pointer">M</div>
                                <div className="hidden md:block">
                                    <h1 className="text-xl font-extrabold text-white uppercase leading-none">Miss Beauty</h1>
                                    <p className="text-[10px] text-yellow-400 font-bold uppercase mt-1">Vietnam 2026</p>
                                </div>
                            </div>

                            <nav className="hidden md:flex space-x-8 h-full items-center">
                                <NavItem tab="overview" label="Tổng quan" icon={HomeIcon} />
                                <NavItem tab="profile" label="Hồ sơ" icon={UserIcon} />
                                <NavItem tab="notifications" label="Hỗ trợ" icon={ChatIcon} />
                                <NavItem tab="gallery" label="Thư viện" icon={PhotoIcon} />
                                <NavItem tab="stats" label="Kết quả" icon={ChartBarIcon} />
                            </nav>

                            <div className="flex items-center space-x-4">
                                <p className="text-sm font-bold text-white leading-tight">{profile?.personalInfo?.fullName}</p>
                                <button onClick={handleLogout} className="p-2.5 rounded-xl text-blue-200 hover:text-white"><LogoutIcon /></button>
                            </div>
                        </div>
                    </div>
                </header>

                <main className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8 mb-20 md:mb-8 w-full">
                    {activeTab === 'overview' && (
                        <div className="animate-fade-in-up">
                            <div className="relative rounded-3xl overflow-hidden bg-white/10 backdrop-blur-md border border-white/20 shadow-2xl mb-12 p-8 md:p-12 text-white">
                                <div className="relative z-10 flex flex-col md:flex-row items-center gap-8">
                                    <div className="relative h-32 w-32 bg-white/20 backdrop-blur-md rounded-full flex items-center justify-center text-5xl font-serif text-white border-2 border-white/50 shadow-inner overflow-hidden">
                                        {profile?.portfolio?.avatarUrl ? (
                                            <img src={getFullImageUrl(profile.portfolio.avatarUrl)} alt="Avatar" className="w-full h-full object-cover" />
                                        ) : (
                                            profile?.personalInfo?.fullName?.charAt(0) || "U"
                                        )}
                                    </div>

                                    <div className="text-center md:text-left flex-1">
                                        <h2 className="text-4xl md:text-5xl font-serif font-bold mb-2">{profile?.personalInfo?.fullName}</h2>
                                        <p className="text-blue-50 text-lg font-light mb-6">SBD: {profile?.sbd || "---"}</p>
                                        <div className="flex gap-2">
                                            {(() => {
                                                const status = profile?.status?.toLowerCase() || "pending";
                                                const config: any = {
                                                    pending: { label: "Đang Chờ Duyệt", bg: "bg-yellow-400", text: "text-blue-900" },
                                                    approved: { label: "Đã Chính Thức", bg: "bg-green-500", text: "text-white" },
                                                    rejected: { label: "Cần Bổ Sung", bg: "bg-red-500", text: "text-white" },
                                                    draft: { label: "Bản Nháp", bg: "bg-gray-400", text: "text-white" }
                                                };
                                                const current = config[status] || config["pending"];
                                                return (
                                                    <span className={`px-4 py-2 ${current.bg} ${current.text} font-bold rounded-full uppercase text-xs shadow-lg`}>
                                                        {current.label}
                                                    </span>
                                                );
                                            })()}
                                        </div>
                                    </div>
                                </div>
                            </div>

                            {/* Announcements Section */}
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                                <div className="bg-white/10 backdrop-blur-md rounded-3xl p-8 border border-white/20 shadow-xl text-white">
                                    <div className="flex items-center gap-3 mb-6">
                                        <div className="p-2 bg-yellow-400 rounded-lg text-blue-900">
                                            <BellIcon />
                                        </div>
                                        <h3 className="text-xl font-serif font-bold">Thông Báo Từ BTC</h3>
                                    </div>
                                    <div className="space-y-4">
                                        {[
                                            { id: 1, date: "07/01/2026", title: "Chào mừng thí sinh Miss Tourism 2026", content: "Hệ thống đã chính thức mở cổng đăng ký. Vui lòng cập nhật hồ sơ đầy đủ trước ngày 15/01." },
                                            { id: 2, date: "05/01/2026", title: "Hướng dẫn chụp ảnh Profile", content: "Quy định về trang phục và kích thước ảnh đã được cập nhật. Vui lòng xem kỹ hướng dẫn." },
                                        ].map((item) => (
                                            <div key={item.id} className="bg-white/5 rounded-xl p-4 border border-white/10 hover:bg-white/10 transition-colors cursor-pointer group">
                                                <div className="flex justify-between items-start mb-2">
                                                    <h4 className="font-bold text-yellow-400 group-hover:text-yellow-300">{item.title}</h4>
                                                    <span className="text-[10px] bg-blue-900/50 px-2 py-1 rounded text-blue-200">{item.date}</span>
                                                </div>
                                                <p className="text-sm text-gray-200 line-clamp-2">{item.content}</p>
                                            </div>
                                        ))}
                                    </div>
                                </div>

                                <div className="bg-white/10 backdrop-blur-md rounded-3xl p-8 border border-white/20 shadow-xl text-white">
                                    <div className="flex items-center gap-3 mb-6">
                                        <div className="p-2 bg-blue-500 rounded-lg text-white">
                                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-6 h-6">
                                                <path strokeLinecap="round" strokeLinejoin="round" d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 012.25-2.25h13.5A2.25 2.25 0 0121 7.5v11.25m-18 0h18M5.25 12h13.5h-13.5zm0 5.25h13.5h-13.5z" />
                                            </svg>
                                        </div>
                                        <h3 className="text-xl font-serif font-bold">Lịch Trình</h3>
                                    </div>
                                    <div className="space-y-4">
                                        <div className="flex gap-4 items-center">
                                            <div className="flex-shrink-0 w-16 text-center bg-white/5 rounded-lg p-2 border border-white/10">
                                                <span className="block text-xl font-bold text-yellow-400">20</span>
                                                <span className="block text-xs uppercase text-blue-200">Tháng 1</span>
                                            </div>
                                            <div>
                                                <h4 className="font-bold text-white">Vòng Sơ Khảo TP.HCM</h4>
                                                <p className="text-xs text-blue-200">08:00 - Trung tâm Hội nghị Gem Center</p>
                                            </div>
                                        </div>
                                        <div className="flex gap-4 items-center opacity-70">
                                            <div className="flex-shrink-0 w-16 text-center bg-white/5 rounded-lg p-2 border border-white/10">
                                                <span className="block text-xl font-bold text-gray-400">25</span>
                                                <span className="block text-xs uppercase text-gray-500">Tháng 1</span>
                                            </div>
                                            <div>
                                                <h4 className="font-bold text-white">Vòng Sơ Khảo Hà Nội</h4>
                                                <p className="text-xs text-blue-200">08:00 - Khách sạn Melia</p>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    )}

                    {activeTab === 'profile' && profile && (() => {
                        const canEdit = (profile.status || 'pending').toLowerCase() === 'pending';
                        return (
                            <div className="animate-fade-in-up bg-white/10 backdrop-blur-md rounded-3xl p-8 border border-white/20 text-white shadow-2xl">
                                <div className="flex justify-between items-center mb-6">
                                    <h3 className="text-2xl font-serif font-bold">Thông Tin Cá nhân</h3>
                                    {canEdit ? (
                                        <span className="text-sm bg-yellow-400 text-blue-900 px-3 py-1 rounded-full font-bold">Cho phép chỉnh sửa</span>
                                    ) : (
                                        <span className="text-sm bg-red-500 text-white px-3 py-1 rounded-full font-bold">Đã khóa (Chỉ xem)</span>
                                    )}
                                </div>

                                <form onSubmit={(e) => {
                                    e.preventDefault();
                                    const formData = new FormData(e.currentTarget);
                                    const input: any = {};
                                    formData.forEach((value, key) => input[key] = value);

                                    // Convert types
                                    input.height = parseFloat(input.height);
                                    input.weight = parseFloat(input.weight);

                                    // Handle Social Links (Split inputs)
                                    const socialLinks = [
                                        input.facebook,
                                        input.instagram,
                                        input.tiktok
                                    ].filter(s => s && s.trim() !== "");
                                    input.socialLinks = socialLinks;
                                    delete input.facebook;
                                    delete input.instagram;
                                    delete input.tiktok;

                                    if (input.galleryUrls) input.galleryUrls = input.galleryUrls.split('\n').map((s: string) => s.trim()).filter((s: string) => s !== "");

                                    const token = localStorage.getItem("token");
                                    apiRequest(UPDATE_PROFILE_MUTATION, { input }, token)
                                        .then((data: any) => {
                                            alert("Cập nhật thành công!");
                                            if (data && data.updateContestantProfile) {
                                                setProfile(data.updateContestantProfile);
                                            }
                                        })
                                        .catch((err: any) => alert("Lỗi: " + err.message));
                                }}>
                                    <div className="space-y-8">
                                        {/* 1. Avatar Section (Top) */}
                                        <div className="flex flex-col items-center justify-center p-6 bg-white/5 rounded-2xl border border-white/10">
                                            <div className="relative group w-32 h-32">
                                                {(() => {
                                                    return (
                                                        <div className="w-full h-full rounded-full overflow-hidden border-4 border-white shadow-xl bg-gray-100">
                                                            {profile.portfolio?.avatarUrl ? (
                                                                <img src={getFullImageUrl(profile.portfolio.avatarUrl)} alt="Avatar" className="w-full h-full object-cover" />
                                                            ) : (
                                                                <div className="w-full h-full flex items-center justify-center text-gray-300 bg-white/10">
                                                                    <UserIcon />
                                                                </div>
                                                            )}
                                                        </div>
                                                    );
                                                })()}

                                                {canEdit && (
                                                    <>
                                                        <div className="absolute inset-0 bg-black/40 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer">
                                                            <span className="text-white text-xs font-bold">Đổi Ảnh</span>
                                                        </div>
                                                        <input
                                                            type="file"
                                                            accept="image/*"
                                                            className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                                                            onChange={(e) => handleFileUpload(e, 'avatarUrl')}
                                                        />
                                                        {uploading && <div className="absolute inset-0 flex items-center justify-center bg-black/50 rounded-full"><div className="animate-spin h-8 w-8 border-2 border-white rounded-full border-t-transparent"></div></div>}
                                                    </>
                                                )}
                                            </div>
                                            <input type="hidden" name="avatarUrl" value={profile.portfolio?.avatarUrl || ''} />
                                            <p className="text-sm text-blue-200 mt-3 font-light">Chạm vào hình để thay đổi ảnh đại diện</p>
                                        </div>

                                        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                                            {/* 2. Physical Info (Priority) */}
                                            <div className="space-y-4">
                                                <h4 className="font-bold text-yellow-400 uppercase text-sm border-b border-white/20 pb-2">Chỉ số hình thể</h4>
                                                <div className="grid grid-cols-2 gap-4">
                                                    <div>
                                                        <label className="block text-xs uppercase text-blue-200 mb-1">Chiều cao (cm)</label>
                                                        <input name="height" type="number" step="0.1" defaultValue={profile.physicalInfo?.height} disabled={!canEdit} className="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-yellow-400 disabled:opacity-50" />
                                                    </div>
                                                    <div>
                                                        <label className="block text-xs uppercase text-blue-200 mb-1">Cân nặng (kg)</label>
                                                        <input name="weight" type="number" step="0.1" defaultValue={profile.physicalInfo?.weight} disabled={!canEdit} className="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-yellow-400 disabled:opacity-50" />
                                                    </div>
                                                </div>
                                                <div>
                                                    <label className="block text-xs uppercase text-blue-200 mb-1">Số đo 3 vòng (cm)</label>
                                                    <input name="measurements" defaultValue={profile.physicalInfo?.measurements} disabled={!canEdit} className="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-yellow-400 disabled:opacity-50" />
                                                </div>
                                            </div>

                                            {/* 3. Personal Info */}
                                            <div className="space-y-4">
                                                <h4 className="font-bold text-yellow-400 uppercase text-sm border-b border-white/20 pb-2">Thông tin cá nhân</h4>
                                                <div>
                                                    <label className="block text-xs uppercase text-blue-200 mb-1">Họ và tên</label>
                                                    <input name="fullName" defaultValue={profile.personalInfo?.fullName} disabled={!canEdit} className="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-yellow-400 disabled:opacity-50" />
                                                </div>
                                                <div className="grid grid-cols-2 gap-4">
                                                    <div>
                                                        <label className="block text-xs uppercase text-blue-200 mb-1">Số điện thoại</label>
                                                        <input name="phone" defaultValue={profile.personalInfo?.phone} disabled={!canEdit} className="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-yellow-400 disabled:opacity-50" />
                                                    </div>
                                                    <div>
                                                        <label className="block text-xs uppercase text-blue-200 mb-1">Email</label>
                                                        <input name="email" defaultValue={profile.personalInfo?.email} disabled={!canEdit} className="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-yellow-400 disabled:opacity-50" />
                                                    </div>
                                                </div>
                                                <div>
                                                    <label className="block text-xs uppercase text-blue-200 mb-1">Địa chỉ</label>
                                                    <input name="address" defaultValue={profile.personalInfo?.address} disabled={!canEdit} className="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-yellow-400 disabled:opacity-50" />
                                                </div>
                                                <div>
                                                    <label className="block text-xs uppercase text-blue-200 mb-1">Nghề nghiệp</label>
                                                    <input name="job" defaultValue={profile.personalInfo?.job} disabled={!canEdit} className="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-yellow-400 disabled:opacity-50" />
                                                </div>
                                            </div>

                                            {/* 4. Social Links (Split) */}
                                            <div className="space-y-4">
                                                <h4 className="font-bold text-yellow-400 uppercase text-sm border-b border-white/20 pb-2">Mạng xã hội</h4>
                                                {(() => {
                                                    const links = profile.portfolio?.socialLinks || [];
                                                    const getLink = (keyword: string) => links.find((l: string) => l.toLowerCase().includes(keyword)) || "";
                                                    return (
                                                        <>
                                                            <div>
                                                                <label className="block text-xs uppercase text-blue-200 mb-1">Facebook</label>
                                                                <input name="facebook" defaultValue={getLink('facebook')} disabled={!canEdit} className="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-yellow-400 disabled:opacity-50" placeholder="https://facebook.com/..." />
                                                            </div>
                                                            <div>
                                                                <label className="block text-xs uppercase text-blue-200 mb-1">Instagram</label>
                                                                <input name="instagram" defaultValue={getLink('instagram')} disabled={!canEdit} className="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-yellow-400 disabled:opacity-50" placeholder="https://instagram.com/..." />
                                                            </div>
                                                            <div>
                                                                <label className="block text-xs uppercase text-blue-200 mb-1">TikTok</label>
                                                                <input name="tiktok" defaultValue={getLink('tiktok')} disabled={!canEdit} className="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-yellow-400 disabled:opacity-50" placeholder="https://tiktok.com/..." />
                                                            </div>
                                                        </>
                                                    );
                                                })()}
                                            </div>

                                            {/* 5. Gallery & Intro */}
                                            <div className="space-y-4">
                                                <h4 className="font-bold text-yellow-400 uppercase text-sm border-b border-white/20 pb-2">Khác</h4>
                                                <div>
                                                    <label className="block text-xs uppercase text-blue-200 mb-1">Thư viện ảnh (URLs)</label>
                                                    <textarea name="galleryUrls" rows={3} defaultValue={profile.portfolio?.galleryUrls?.join('\n')} disabled={!canEdit} className="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-yellow-400 disabled:opacity-50" placeholder="Mỗi dòng một link ảnh..." />
                                                    <p className="text-[10px] text-blue-300 mt-1">Nhập mỗi link trên một dòng</p>
                                                </div>
                                                <div>
                                                    <label className="block text-xs uppercase text-blue-200 mb-1">Giới thiệu bản thân</label>
                                                    <textarea name="introduction" rows={4} defaultValue={profile.portfolio?.introduction} disabled={!canEdit} className="w-full bg-white/5 border border-white/10 rounded-lg px-4 py-2 focus:outline-none focus:border-yellow-400 disabled:opacity-50" />
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                    {canEdit && (
                                        <div className="mt-8 flex justify-end">
                                            <button type="submit" className="bg-yellow-400 text-blue-900 font-bold px-8 py-3 rounded-xl hover:bg-yellow-300 transition-all shadow-lg hover:shadow-yellow-400/20">
                                                Lưu Thay Đổi
                                            </button>
                                        </div>
                                    )}
                                </form>
                            </div>
                        );
                    })()}

                    {activeTab === 'notifications' && (
                        <div className="animate-fade-in-up">
                            <div className="flex justify-between items-center mb-6">
                                <h3 className="text-2xl font-serif font-bold text-white">Hộp thư hỗ trợ</h3>
                                <button onClick={() => setIsFeedbackModalOpen(true)} className="bg-yellow-400 hover:bg-yellow-300 text-blue-900 px-4 py-2 rounded-lg font-bold shadow-lg transition-all flex items-center gap-2">
                                    <ChatIcon /> Gửi yêu cầu mới
                                </button>
                            </div>

                            {loadingFeedbacks ? (
                                <div className="text-center text-white py-12">Đang tải dữ liệu...</div>
                            ) : feedbacks.length === 0 ? (
                                <div className="text-center text-blue-200 py-12 bg-white/5 rounded-3xl border border-white/10">
                                    <p className="text-lg">Bạn chưa gửi yêu cầu hỗ trợ nào.</p>
                                    <p className="text-sm mt-2">Nếu cần giúp đỡ, hãy nhấn nút "Gửi yêu cầu mới".</p>
                                </div>
                            ) : (
                                <div className="space-y-4">
                                    {feedbacks.map((item: any) => (
                                        <div key={item.id} className="bg-white/10 backdrop-blur-md rounded-2xl p-6 border border-white/20 shadow-lg text-white">
                                            <div className="flex flex-col md:flex-row justify-between gap-4 mb-4 border-b border-white/10 pb-4">
                                                <div>
                                                    <div className="flex items-center gap-3">
                                                        {item.status === 'reply' ? (
                                                            <span className="bg-green-500 text-white text-[10px] font-bold px-2 py-1 rounded uppercase">Đã phản hồi</span>
                                                        ) : (
                                                            <span className="bg-yellow-500/50 text-yellow-100 text-[10px] font-bold px-2 py-1 rounded uppercase">Đang xử lý</span>
                                                        )}
                                                        <h4 className="font-bold text-lg text-yellow-50">{item.title}</h4>
                                                    </div>
                                                    <p className="text-xs text-blue-200 mt-1">Gửi lúc: {new Date(item.createdAt).toLocaleString('vi-VN')}</p>
                                                </div>
                                                <div>
                                                    <span className="text-xs uppercase tracking-wider font-bold opacity-70 border border-white/30 px-2 py-1 rounded">
                                                        {item.type === 'complaint' ? 'Khiếu nại' : item.type === 'proposal' ? 'Đề xuất' : 'Hỗ trợ'}
                                                    </span>
                                                </div>
                                            </div>

                                            <div className="space-y-4">
                                                <div className="bg-black/20 p-4 rounded-xl">
                                                    <p className="font-bold text-xs text-blue-300 mb-1">Nội dung gửi:</p>
                                                    <p className="text-sm whitespace-pre-wrap">{item.content}</p>
                                                </div>

                                                {item.reply && (
                                                    <div className="bg-green-900/40 border border-green-500/30 p-4 rounded-xl ml-8 relative">
                                                        <div className="absolute -left-2 top-4 w-4 h-4 bg-green-900/40 border-l border-b border-green-500/30 transform rotate-45"></div>
                                                        <p className="font-bold text-xs text-green-300 mb-1">Phản hồi từ BTC:</p>
                                                        <p className="text-sm text-white whitespace-pre-wrap">{item.reply}</p>
                                                    </div>
                                                )}
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                    )}

                    {activeTab === 'gallery' && (
                        <div className="text-center text-white">Gallery Feature Coming Soon (Read Only)</div>
                    )}

                    {activeTab === 'stats' && (
                        <div className="animate-fade-in-up">
                            <h3 className="text-2xl font-serif font-bold text-white mb-6">Kết quả thi</h3>

                            {loadingScores ? (
                                <div className="text-center text-white py-12">Đang tải dữ liệu...</div>
                            ) : scores.length === 0 ? (
                                <div className="text-center text-blue-200 py-12 bg-white/5 rounded-3xl border border-white/10">
                                    <p className="text-lg">Chưa có kết quả thi.</p>
                                    <p className="text-sm mt-2">Vui lòng quay lại sau khi Ban Giám Khảo công bố điểm.</p>
                                </div>
                            ) : (
                                <div className="space-y-6">
                                    {scores.map((score: any) => (
                                        <div key={score.id} className="bg-white/10 backdrop-blur-md rounded-3xl p-8 border border-white/20 shadow-2xl text-white">
                                            <div className="flex justify-between items-center mb-6 border-b border-white/10 pb-4">
                                                <div>
                                                    <h4 className="text-xl font-bold text-yellow-400 uppercase tracking-widest">{score.roundId}</h4>
                                                    <p className="text-xs text-blue-200 mt-1">SBD: {score.sbd}</p>
                                                </div>
                                                <div className="text-center">
                                                    <div className="text-4xl font-serif font-bold text-white">{score.totalScore.toFixed(1)}</div>
                                                    <div className="text-[10px] uppercase text-blue-200">Tổng điểm</div>
                                                </div>
                                            </div>

                                            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                                                {/* Criteria Breakdown */}
                                                <div>
                                                    <h5 className="font-bold text-sm text-blue-200 uppercase mb-4 border-b border-white/10 pb-2">Chi tiết điểm thành phần</h5>
                                                    <div className="space-y-3">
                                                        {Object.entries(score.criteriaScores || {}).map(([key, value]: [string, any]) => (
                                                            <div key={key} className="flex justify-between items-center bg-white/5 px-4 py-2 rounded-lg">
                                                                <span className="text-sm font-medium text-gray-100 capitalize">{key.replace(/_/g, ' ')}</span>
                                                                <span className="font-bold text-yellow-400">{Number(value).toFixed(1)}</span>
                                                            </div>
                                                        ))}
                                                    </div>
                                                </div>

                                                {/* Judge Comments */}
                                                <div>
                                                    <h5 className="font-bold text-sm text-blue-200 uppercase mb-4 border-b border-white/10 pb-2">Nhận xét từ Ban Giám Khảo</h5>
                                                    {score.comment ? (
                                                        <div className="bg-gradient-to-br from-blue-900/40 to-black/20 p-6 rounded-2xl border border-white/10 relative">
                                                            <svg className="absolute top-4 left-4 w-6 h-6 text-white/10" viewBox="0 0 24 24" fill="currentColor"><path d="M14.017 21L14.017 18C14.017 16.8954 13.1216 16 12.017 16H9C9 14.9391 9.39063 13.9189 10.0933 13.1553C10.8252 12.3599 11.8385 11.8741 12.9175 11.8105L14.017 11.724L14.017 9.17647L12.7231 9.25547C10.9168 9.36629 9.17482 10.1504 7.91508 11.5192C6.67138 12.8706 5.99997 14.6534 5.99997 16.5V21H14.017ZM21 21L21 18C21 16.8954 20.1046 16 19 16H15.983C15.983 14.9391 16.3736 13.9189 17.0763 13.1553C17.8082 12.3599 18.8215 11.8741 19.9005 11.8105L21 11.724L21 9.17647L19.7061 9.25547C17.8998 9.36629 16.1578 10.1504 14.8981 11.5192C13.6544 12.8706 12.983 14.6534 12.983 16.5V21H21Z" /></svg>
                                                            <p className="text-sm text-gray-200 italic leading-relaxed pl-2 relative z-10">"{score.comment}"</p>
                                                        </div>
                                                    ) : (
                                                        <p className="text-sm text-gray-400 italic">Không có nhận xét nào.</p>
                                                    )}
                                                </div>
                                            </div>

                                            <div className="mt-6 pt-4 border-t border-white/10 flex justify-between items-center opacity-60">
                                                <span className="text-[10px] uppercase">ID: {score.id}</span>
                                                <span className="text-[10px] uppercase">{new Date(score.createdAt).toLocaleDateString('vi-VN')}</span>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                    )}
                </main>

                {/* Floating Action Button */}
                <button
                    onClick={() => setIsFeedbackModalOpen(true)}
                    className="fixed bottom-8 right-8 bg-gradient-to-r from-pink-500 to-rose-500 text-white p-4 rounded-full shadow-[0_4px_20px_rgba(236,72,153,0.5)] hover:scale-110 hover:shadow-[0_4px_25px_rgba(236,72,153,0.7)] transition-all z-40 group"
                    title="Gửi phản hồi"
                >
                    <ChatIcon />
                    <span className="absolute right-full mr-3 top-1/2 -translate-y-1/2 bg-white text-blue-900 px-3 py-1 rounded-lg text-sm font-bold opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap shadow-lg pointer-events-none">
                        Gửi hỗ trợ
                    </span>
                </button>

                {/* Feedback Modal */}
                {isFeedbackModalOpen && (
                    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
                        <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" onClick={() => setIsFeedbackModalOpen(false)}></div>
                        <div className="relative bg-white rounded-2xl shadow-2xl w-full max-w-lg overflow-hidden animate-fade-in-up">
                            <div className="bg-gradient-to-r from-blue-900 to-blue-800 p-6 text-white flex justify-between items-center">
                                <h3 className="font-serif font-bold text-xl flex items-center gap-2">
                                    <ChatIcon /> Gửi yêu cầu hỗ trợ
                                </h3>
                                <button onClick={() => setIsFeedbackModalOpen(false)} className="text-white/70 hover:text-white text-2xl leading-none">&times;</button>
                            </div>
                            <form onSubmit={handleSendFeedback} className="p-6 space-y-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Loại yêu cầu</label>
                                    <select
                                        className="w-full p-3 border border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 outline-none"
                                        value={feedbackType}
                                        onChange={(e) => setFeedbackType(e.target.value)}
                                    >
                                        <option value="support">Hỗ trợ kỹ thuật</option>
                                        <option value="complaint">Khiếu nại kết quả/quy trình</option>
                                        <option value="proposal">Đề xuất ý kiến</option>
                                    </select>
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Tiêu đề (Ngắn gọn)</label>
                                    <input
                                        className="w-full p-3 border border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 outline-none"
                                        placeholder="Ví dụ: Sai thông tin ngày sinh..."
                                        value={feedbackTitle}
                                        onChange={(e) => setFeedbackTitle(e.target.value)}
                                        required
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Nội dung chi tiết</label>
                                    <textarea
                                        className="w-full p-3 border border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 outline-none"
                                        rows={5}
                                        placeholder="Mô tả chi tiết vấn đề của bạn..."
                                        value={feedbackContent}
                                        onChange={(e) => setFeedbackContent(e.target.value)}
                                        required
                                    ></textarea>
                                </div>
                                <div className="flex justify-end pt-2">
                                    <button type="button" onClick={() => setIsFeedbackModalOpen(false)} className="px-4 py-2 text-gray-600 hover:text-gray-900 font-medium mr-2">Hủy bỏ</button>
                                    <button type="submit" className="bg-blue-900 hover:bg-blue-800 text-white px-6 py-2 rounded-xl font-bold shadow-lg transition-transform hover:-translate-y-0.5">
                                        Gửi ngay
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};
