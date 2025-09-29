"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAppSelector } from "@/hooks/redux";
import ProtectedRoute from "@/components/auth/protected-route";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { apiClient, Booking } from "@/lib/api-client";
import { toast } from "sonner";
import { Calendar, Car, DollarSign, Clock, ArrowLeft } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";

export default function BookingsPage() {
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [loading, setLoading] = useState(true);
  const { user } = useAppSelector((state) => state.auth);
  const router = useRouter();

  useEffect(() => {
    const fetchBookings = async () => {
      if (!user?.id) return;

      try {
        const data = await apiClient.getBookingsByCustomer(user.id);
        setBookings(data || []); // Handle null/undefined response
      } catch (error) {
        toast.error("Failed to fetch bookings");
        console.error("Error fetching bookings:", error);
        setBookings([]); // Set empty array on error
      } finally {
        setLoading(false);
      }
    };

    fetchBookings();
  }, [user?.id]);

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case "confirmed":
        return "bg-green-100 text-green-800";
      case "pending":
        return "bg-yellow-100 text-yellow-800";
      case "cancelled":
        return "bg-red-100 text-red-800";
      case "completed":
        return "bg-blue-100 text-blue-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  const calculateDays = (startDate: string, endDate: string) => {
    const start = new Date(startDate);
    const end = new Date(endDate);
    const diffTime = Math.abs(end.getTime() - start.getTime());
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    return diffDays;
  };

  if (loading) {
    return (
      <ProtectedRoute>
        <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 p-8">
          <div className="max-w-6xl mx-auto">
            <div className="flex items-center gap-4 mb-8">
              <Skeleton className="h-10 w-10" />
              <Skeleton className="h-8 w-48" />
            </div>
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
              {[...Array(6)].map((_, i) => (
                <Card key={i} className="overflow-hidden">
                  <CardHeader>
                    <Skeleton className="h-6 w-3/4" />
                    <Skeleton className="h-4 w-1/2" />
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-3">
                      <Skeleton className="h-4 w-full" />
                      <Skeleton className="h-4 w-full" />
                      <Skeleton className="h-4 w-3/4" />
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>
        </div>
      </ProtectedRoute>
    );
  }

  return (
    <ProtectedRoute>
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 p-8">
        <div className="max-w-6xl mx-auto">
          <div className="flex items-center gap-4 mb-8">
            <Button
              variant="outline"
              size="icon"
              onClick={() => router.push("/dashboard")}
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">My Bookings</h1>
              <p className="text-gray-600">Manage your car rental bookings</p>
            </div>
          </div>

          {bookings.length === 0 ? (
            <Card className="text-center py-12">
              <CardContent>
                <Car className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <h3 className="text-lg font-semibold text-gray-900 mb-2">
                  No bookings yet
                </h3>
                <p className="text-gray-600 mb-4">
                  You haven&apos;t made any car bookings yet. Start exploring
                  our cars!
                </p>
                <Button onClick={() => router.push("/cars")}>
                  Browse Cars
                </Button>
              </CardContent>
            </Card>
          ) : (
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
              {bookings.map((booking) => (
                <Card
                  key={booking.id}
                  className="overflow-hidden hover:shadow-lg transition-shadow"
                >
                  <CardHeader>
                    <div className="flex justify-between items-start">
                      <div>
                        <CardTitle className="text-lg">
                          {booking.car?.brand} {booking.car?.model}
                        </CardTitle>
                        <CardDescription>Booking #{booking.id}</CardDescription>
                      </div>
                      <Badge className={getStatusColor(booking.status)}>
                        {booking.status}
                      </Badge>
                    </div>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="flex items-center gap-2 text-sm text-gray-600">
                      <Calendar className="h-4 w-4" />
                      <span>
                        {formatDate(booking.start_date)} -{" "}
                        {formatDate(booking.end_date)}
                      </span>
                    </div>

                    <div className="flex items-center gap-2 text-sm text-gray-600">
                      <Clock className="h-4 w-4" />
                      <span>
                        {calculateDays(booking.start_date, booking.end_date)}{" "}
                        days
                      </span>
                    </div>

                    <div className="flex items-center gap-2 text-sm text-gray-600">
                      <DollarSign className="h-4 w-4" />
                      <span className="font-semibold text-green-600">
                        ${booking.total_cost.toLocaleString()}
                      </span>
                    </div>

                    {booking.car && (
                      <div className="pt-2 border-t">
                        <div className="flex justify-between text-sm text-gray-600">
                          <span>Year:</span>
                          <span>{booking.car.year}</span>
                        </div>
                        <div className="flex justify-between text-sm text-gray-600">
                          <span>Fuel:</span>
                          <span>{booking.car.fuel_type}</span>
                        </div>
                        <div className="flex justify-between text-sm text-gray-600">
                          <span>Transmission:</span>
                          <span>
                            {booking.car.engine?.transmission || "N/A"}
                          </span>
                        </div>
                      </div>
                    )}

                    <div className="flex gap-2 pt-2">
                      <Button
                        variant="outline"
                        size="sm"
                        className="flex-1"
                        onClick={() => router.push(`/cars/${booking.car_id}`)}
                      >
                        View Car
                      </Button>
                      {booking.status === "pending" && (
                        <Button
                          variant="destructive"
                          size="sm"
                          onClick={async () => {
                            try {
                              await apiClient.deleteBooking(booking.id);
                              setBookings(
                                bookings.filter((b) => b.id !== booking.id)
                              );
                              toast.success("Booking cancelled successfully");
                            } catch {
                              toast.error("Failed to cancel booking");
                            }
                          }}
                        >
                          Cancel
                        </Button>
                      )}
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </div>
      </div>
    </ProtectedRoute>
  );
}
