import { useNavigate } from "@tanstack/react-router";
import { useEffect } from "react";
import { apiRequest } from "../lib/api";
import { useAuthStore } from "../store/useAuthStore";

const MY_PROFILE_QUERY_CHECK = `
query MyProfile {
  myProfile {
    id
  }
}
`;

export const ContestantDispatcher = () => {
    const navigate = useNavigate();
    const { token, isAuthenticated } = useAuthStore();

    useEffect(() => {
        if (!isAuthenticated || !token) {
            navigate({ to: "/login" as "/login" });
            return;
        }

        apiRequest(MY_PROFILE_QUERY_CHECK, {}, token)
            .then((data: any) => {
                if (data && data.myProfile) {
                    // Registered -> Dashboard
                    navigate({ to: "/contestant/dashboard" as "/contestant/dashboard" });
                } else {
                    // Not Registered -> Register Form
                    navigate({ to: "/contestant/register" as "/contestant/register" });
                }
            })
            .catch(() => {
                // If error, assume needs login or something wrong, sending to login is safest
                navigate({ to: "/login" as "/login" });
            });
    }, [navigate, isAuthenticated, token]);

    return (
        <div className="min-h-screen flex items-center justify-center bg-white">
            <div className="flex flex-col items-center animate-pulse">
                <div className="h-12 w-12 bg-blue-900 rounded-full mb-4"></div>
                <div className="text-gray-400 font-light tracking-widest uppercase">Checking Profile...</div>
            </div>
        </div>
    );
};
