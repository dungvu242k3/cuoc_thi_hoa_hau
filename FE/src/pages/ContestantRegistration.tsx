import { useNavigate } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { apiRequest } from "../lib/api";

const UserIcon = () => (<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5"><path strokeLinecap="round" strokeLinejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" /></svg>);

const MY_PROFILE_QUERY = `
query MyProfile {
  myProfile {
    id
    personalInfo {
      fullName
      dob
      nationality
      identityCard
      phone
      email
      address
      job
    }
    physicalInfo {
      height
      weight
      measurements
    }
    skillEducation {
        educationLevel
        languages
        skills
    }
    portfolio {
        introduction
        socialLinks
        avatarUrl
    }
  }
}
`;

const CREATE_PROFILE_MUTATION = `
mutation CreateContestantProfile($input: CreateContestantInput!) {
    createContestantProfile(input: $input) {
        id
        status
        personalInfo { fullName }
    }
}
`;

export const ContestantRegistration = () => {
    const navigate = useNavigate();
    const { register, handleSubmit, setValue, watch } = useForm();
    const [isSaving, setIsSaving] = useState(false);
    const [uploading, setUploading] = useState(false);
    const [avatarUrl, setAvatarUrl] = useState("");

    // Monitor avatar field for preview
    const watchedAvatar = watch("avatarUrl");
    useEffect(() => {
        if (watchedAvatar) setAvatarUrl(watchedAvatar);
    }, [watchedAvatar]);

    useEffect(() => {
        const token = localStorage.getItem("token");
        if (!token) {
            navigate({ to: "/login" });
            return;
        }

        // Check if already registered
        apiRequest(MY_PROFILE_QUERY, {}, token)
            .then((data: any) => {
                if (data && data.myProfile) {
                    // Already registered -> Go to Dashboard
                    navigate({ to: "/contestant/dashboard" });
                } else {
                    // Pre-fill email/phone if available from other sources? 
                    // For now just let them fill.
                }
            })
            .catch(() => { });
    }, [navigate]);

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
            setValue(fieldName, data.url);
        } catch (error: any) {
            console.error("Upload error:", error);
            alert("Lỗi tải ảnh lên.");
        } finally {
            setUploading(false);
        }
    };

    const onSubmit = async (data: any) => {
        setIsSaving(true);
        const token = localStorage.getItem("token");

        const input = {
            fullName: data.fullName,
            dob: data.dob,
            phone: data.phone,
            email: data.email,
            address: data.address,
            job: data.job,
            nationality: data.nationality,
            height: isNaN(parseFloat(data.height)) ? 0 : parseFloat(data.height),
            weight: isNaN(parseFloat(data.weight)) ? 0 : parseFloat(data.weight),
            measurements: data.measurements,
            educationLevel: data.educationLevel,
            languages: data.languages ? data.languages.split(',').map((s: string) => s.trim()) : [],
            skills: data.skills ? data.skills.split(',').map((s: string) => s.trim()) : [],
            introduction: data.introduction,
            socialLinks: [data.facebook, data.instagram].filter(Boolean),
            avatarUrl: data.avatarUrl,
        };

        try {
            await apiRequest(CREATE_PROFILE_MUTATION, { input }, token);
            alert("Đăng ký hồ sơ thành công!");
            navigate({ to: "/contestant/dashboard" });
            window.location.reload(); // Ensure dashboard loads fresh data
        } catch (error: any) {
            console.error("Registration Error:", error);
            const msg = error.message || JSON.stringify(error);
            alert(`Có lỗi xảy ra: ${msg}`);
        } finally {
            setIsSaving(false);
        }
    };

    return (
        <div className="min-h-screen bg-gradient-to-br from-blue-50 to-white py-12 px-4 sm:px-6 lg:px-8 font-sans text-slate-800">
            <div className="max-w-4xl mx-auto">
                <div className="text-center mb-10">
                    <h2 className="text-3xl font-serif font-bold text-blue-900 mb-2">Đăng Ký Dự Thi</h2>
                    <p className="text-slate-500">Miss Tourism Vietnam 2026</p>
                </div>

                <div className="bg-white rounded-2xl shadow-xl border border-gray-100 overflow-hidden">
                    <div className="p-6 border-b border-gray-100 bg-blue-50/50">
                        <h3 className="font-bold text-xl text-blue-950">Thông tin hồ sơ</h3>
                        <p className="text-sm text-slate-500">Vui lòng điền chính xác thông tin</p>
                    </div>

                    <form onSubmit={handleSubmit(onSubmit)} className="p-8 space-y-8">
                        {/* Avatar Upload */}
                        <div className="flex justify-center">
                            <div className="relative group cursor-pointer w-32 h-32">
                                {(() => {
                                    const getFullImageUrl = (path: string) => {
                                        if (!path) return "";
                                        if (path.startsWith("http") || path.startsWith("data:")) return path;
                                        const baseUrl = (import.meta.env.VITE_API_URL || "http://localhost:8080").replace(/\/query\/?$/, "");
                                        return `${baseUrl}${path}`;
                                    };

                                    return (
                                        <div className="w-full h-full rounded-full overflow-hidden border-4 border-white shadow-lg bg-gray-100">
                                            {avatarUrl ? (
                                                <img src={getFullImageUrl(avatarUrl)} alt="Avatar" className="w-full h-full object-cover" />
                                            ) : (
                                                <div className="w-full h-full flex items-center justify-center text-4xl text-gray-300">
                                                    <UserIcon />
                                                </div>
                                            )}
                                        </div>
                                    );
                                })()}
                                <div className="absolute inset-0 bg-black/40 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                                    <span className="text-white text-xs font-bold">Change</span>
                                </div>
                                <input
                                    type="file"
                                    accept="image/*"
                                    className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                                    onChange={(e) => handleFileUpload(e, 'avatarUrl')}
                                />
                                {uploading && <div className="absolute inset-0 flex items-center justify-center bg-white/80 rounded-full"><div className="animate-spin h-6 w-6 border-2 border-blue-600 rounded-full border-t-transparent"></div></div>}
                            </div>
                        </div>

                        {/* Section 1: Personal Info */}
                        <div>
                            <h4 className="text-xs font-bold text-blue-900 uppercase tracking-widest mb-4 border-l-4 border-yellow-400 pl-3">1. Thông tin cá nhân</h4>
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                <div className="md:col-span-2">
                                    <label className="block text-sm font-medium text-slate-700 mb-1">Họ và tên</label>
                                    <input {...register("fullName")} className="w-full p-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none" placeholder="Nguyễn Văn A" required />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-slate-700 mb-1">Ngày sinh</label>
                                    <input type="date" {...register("dob")} className="w-full p-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none" required />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-slate-700 mb-1">Quốc tịch</label>
                                    <input {...register("nationality")} className="w-full p-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none" placeholder="Việt Nam" />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-slate-700 mb-1">Số điện thoại</label>
                                    <input {...register("phone")} className="w-full p-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none" placeholder="09xxxxxxxx" required />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-slate-700 mb-1">Email</label>
                                    <input {...register("email")} className="w-full p-3 rounded-lg border border-gray-200 bg-gray-50" placeholder="Email liên hệ" required />
                                </div>
                                <div className="md:col-span-2">
                                    <label className="block text-sm font-medium text-slate-700 mb-1">Địa chỉ thường trú</label>
                                    <input {...register("address")} className="w-full p-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none" placeholder="Số nhà, đường..." />
                                </div>
                                <div className="md:col-span-2">
                                    <label className="block text-sm font-medium text-slate-700 mb-1">Nghề nghiệp hiện tại</label>
                                    <input {...register("job")} className="w-full p-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none" placeholder="Sinh viên..." />
                                </div>
                            </div>
                        </div>

                        {/* Section 2: Physical Info */}
                        <div className="border-t border-dashed border-gray-200 pt-8">
                            <h4 className="text-xs font-bold text-blue-900 uppercase tracking-widest mb-4 border-l-4 border-yellow-400 pl-3">2. Chỉ số hình thể</h4>
                            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                                <div>
                                    <label className="block text-sm font-medium text-slate-700 mb-1">Chiều cao (cm)</label>
                                    <input type="number" step="0.1" {...register("height")} className="w-full p-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none" placeholder="170" required />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-slate-700 mb-1">Cân nặng (kg)</label>
                                    <input type="number" step="0.1" {...register("weight")} className="w-full p-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none" placeholder="50" required />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-slate-700 mb-1">Số đo 3 vòng (cm)</label>
                                    <input {...register("measurements")} className="w-full p-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none" placeholder="90-60-90" required />
                                </div>
                            </div>
                        </div>

                        {/* Section 3: Education & Skills */}
                        <div className="border-t border-dashed border-gray-200 pt-8">
                            <h4 className="text-xs font-bold text-blue-900 uppercase tracking-widest mb-4 border-l-4 border-yellow-400 pl-3">3. Học vấn & Kỹ năng</h4>
                            <div className="grid grid-cols-1 gap-6">
                                <div>
                                    <label className="block text-sm font-medium text-slate-700 mb-1">Trình độ học vấn</label>
                                    <select {...register("educationLevel")} className="w-full p-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none">
                                        <option value="">-- Chọn trình độ --</option>
                                        <option value="HighSchool">Trung học phổ thông</option>
                                        <option value="College">Cao đẳng</option>
                                        <option value="University">Đại học</option>
                                        <option value="PostGraduate">Sau đại học</option>
                                    </select>
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-slate-700 mb-1">Ngoại ngữ</label>
                                    <input {...register("languages")} className="w-full p-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none" placeholder="Tiếng Anh..." />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-slate-700 mb-1">Năng khiếu / Kỹ năng khác</label>
                                    <input {...register("skills")} className="w-full p-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none" placeholder="Múa, Hát..." />
                                </div>
                            </div>
                        </div>

                        {/* Section 4: Intro & Socials */}
                        <div className="border-t border-dashed border-gray-200 pt-8">
                            <h4 className="text-xs font-bold text-blue-900 uppercase tracking-widest mb-4 border-l-4 border-yellow-400 pl-3">4. Giới thiệu bản thân</h4>
                            <div className="space-y-6">
                                <div>
                                    <label className="block text-sm font-medium text-slate-700 mb-1">Lời giới thiệu</label>
                                    <textarea {...register("introduction")} rows={4} className="w-full p-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none" placeholder="Giới thiệu về bản thân..."></textarea>
                                </div>
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                    <div>
                                        <label className="block text-sm font-medium text-slate-700 mb-1">Link Facebook</label>
                                        <input {...register("facebook")} className="w-full p-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none" placeholder="https://facebook.com/..." />
                                    </div>
                                    <div>
                                        <label className="block text-sm font-medium text-slate-700 mb-1">Link Instagram</label>
                                        <input {...register("instagram")} className="w-full p-3 rounded-lg border border-gray-200 focus:ring-2 focus:ring-blue-500 outline-none" placeholder="https://instagram.com/..." />
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div className="flex justify-end pt-4 border-t border-gray-100">
                            <button
                                type="submit"
                                disabled={isSaving}
                                className="px-8 py-3 bg-gradient-to-r from-blue-900 to-blue-800 text-white font-bold uppercase tracking-wider rounded-xl shadow-lg hover:shadow-blue-900/40 hover:-translate-y-1 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                            >
                                {isSaving ? 'Đang xử lý...' : 'Gửi Hồ Sơ Dự Thi'}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};
