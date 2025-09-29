"use client";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useAppDispatch, useAppSelector } from "@/hooks/redux";
import { logoutUser } from "@/store/slices/authSlice";
import ProtectedRoute from "@/components/auth/protected-route";
import { toast } from "sonner";
import { useRouter } from "next/navigation";

export default function DashboardPage() {
  const { user } = useAppSelector((state) => state.auth);
  const dispatch = useAppDispatch();
  const router = useRouter();

  const handleLogout = async () => {
    try {
      await dispatch(logoutUser());
      toast.success("Logged out successfully");
      router.push("/");
    } catch {
      toast.error("Logout failed");
    }
  };

  return (
    <ProtectedRoute>
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 p-8">
        <div className="max-w-4xl mx-auto">
          <div className="flex justify-between items-center mb-8">
            <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
            <Button onClick={handleLogout} variant="outline">
              Logout
            </Button>
          </div>

          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            <Card>
              <CardHeader>
                <CardTitle>Welcome Back!</CardTitle>
                <CardDescription>Good to see you again</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <p>
                    <span className="font-semibold">Username:</span>{" "}
                    {user?.username}
                  </p>
                  <p>
                    <span className="font-semibold">Email:</span> {user?.email}
                  </p>
                  <p>
                    <span className="font-semibold">Phone:</span> {user?.phone}
                  </p>
                  <p>
                    <span className="font-semibold">Role:</span> {user?.role}
                  </p>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>My Bookings</CardTitle>
                <CardDescription>Manage your car rentals</CardDescription>
              </CardHeader>
              <CardContent>
                <Button
                  className="w-full"
                  onClick={() => router.push("/bookings")}
                >
                  View Bookings
                </Button>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Browse Cars</CardTitle>
                <CardDescription>Find your perfect ride</CardDescription>
              </CardHeader>
              <CardContent>
                <Button className="w-full" onClick={() => router.push("/cars")}>
                  Browse Cars
                </Button>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </ProtectedRoute>
  );
}
