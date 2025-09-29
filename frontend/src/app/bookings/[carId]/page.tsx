"use client";

import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import { useAppSelector } from "@/hooks/redux";
import ProtectedRoute from "@/components/auth/protected-route";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { apiClient, Car } from "@/lib/api-client";
import { toast } from "sonner";
import {
  ArrowLeft,
  Car as CarIcon,
  Calendar,
  User,
  CreditCard,
  MapPin,
  Info,
  Calculator,
} from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import Image from "next/image";

export default function BookingFormPage() {
  const [car, setCar] = useState<Car | null>(null);
  const [loading, setLoading] = useState(true);
  const [bookingLoading, setBookingLoading] = useState(false);

  // Form state
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");
  const [bookingType, setBookingType] = useState("rental");
  const [notes, setNotes] = useState("");

  const { user } = useAppSelector((state) => state.auth);
  const router = useRouter();
  const params = useParams();
  const carId = params.carId as string;

  useEffect(() => {
    const fetchCar = async () => {
      try {
        const data = await apiClient.getCarById(carId);
        setCar(data);
      } catch (error) {
        toast.error("Failed to fetch car details");
        console.error("Error fetching car:", error);
        router.push("/cars");
      } finally {
        setLoading(false);
      }
    };

    if (carId) {
      fetchCar();
    }
  }, [carId, router]);

  const calculateTotalCost = () => {
    if (!startDate || !endDate || !car) return 0;

    const start = new Date(startDate);
    const end = new Date(endDate);
    const diffTime = Math.abs(end.getTime() - start.getTime());
    const diffDays = Math.max(1, Math.ceil(diffTime / (1000 * 60 * 60 * 24)));

    if (bookingType === "rental") {
      return diffDays * car.price.rental_price_daily;
    }
    return car.price.sale_price || car.price.rental_price_daily * diffDays;
  };

  const getDuration = () => {
    if (!startDate || !endDate) return 0;
    const start = new Date(startDate);
    const end = new Date(endDate);
    const diffTime = Math.abs(end.getTime() - start.getTime());
    return Math.max(1, Math.ceil(diffTime / (1000 * 60 * 60 * 24)));
  };

  const handleBooking = async () => {
    if (!startDate || !endDate) {
      toast.error("Please select start and end dates");
      return;
    }

    if (!user?.id) {
      toast.error("Please login to book a car");
      return;
    }

    if (!car?.owner?.id) {
      toast.error("Car owner information not found");
      return;
    }

    setBookingLoading(true);
    try {
      const bookingRequest: {
        customer_id: string;
        car_id: string;
        owner_id: string;
        booking_type: "rental" | "purchase";
        start_date?: string;
        end_date?: string;
        notes?: string;
      } = {
        customer_id: user.id,
        car_id: carId,
        owner_id: car.owner.id,
        booking_type: bookingType as "rental" | "purchase",
      };

      if (bookingType === "rental" && startDate && endDate) {
        bookingRequest.start_date = startDate;
        bookingRequest.end_date = endDate;
      }

      if (notes.trim()) {
        bookingRequest.notes = notes.trim();
      }

      await apiClient.createBooking(bookingRequest);
      toast.success("Booking created successfully!");
      router.push("/bookings");
    } catch (error) {
      toast.error("Failed to create booking");
      console.error("Error creating booking:", error);
    } finally {
      setBookingLoading(false);
    }
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      minimumFractionDigits: 0,
    }).format(price);
  };

  if (loading) {
    return (
      <ProtectedRoute>
        <div className="min-h-screen bg-gradient-to-br from-gray-50 to-white py-8">
          <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="space-y-6">
              <Skeleton className="h-8 w-64" />
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                <Skeleton className="h-96" />
                <Skeleton className="h-96" />
              </div>
            </div>
          </div>
        </div>
      </ProtectedRoute>
    );
  }

  if (!car) {
    return (
      <ProtectedRoute>
        <div className="min-h-screen bg-gradient-to-br from-gray-50 to-white py-8">
          <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
            <h1 className="text-2xl font-bold text-gray-900">Car not found</h1>
            <Button onClick={() => router.push("/cars")} className="mt-4">
              Back to Cars
            </Button>
          </div>
        </div>
      </ProtectedRoute>
    );
  }

  if (!car.is_available) {
    return (
      <ProtectedRoute>
        <div className="min-h-screen bg-gradient-to-br from-gray-50 to-white py-8">
          <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
            <h1 className="text-2xl font-bold text-gray-900">
              Car Not Available
            </h1>
            <p className="text-gray-600 mt-2">
              This car is currently not available for booking.
            </p>
            <div className="flex gap-4 justify-center mt-6">
              <Button
                variant="outline"
                onClick={() => router.push(`/cars/${carId}`)}
              >
                Back to Car Details
              </Button>
              <Button onClick={() => router.push("/cars")}>Browse Cars</Button>
            </div>
          </div>
        </div>
      </ProtectedRoute>
    );
  }

  return (
    <ProtectedRoute>
      <div className="min-h-screen bg-gradient-to-br from-gray-50 to-white py-8">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
          {/* Header */}
          <div className="flex items-center justify-between mb-8">
            <div className="flex items-center gap-4">
              <Button
                variant="outline"
                size="sm"
                onClick={() => router.push(`/cars/${carId}`)}
                className="flex items-center gap-2"
              >
                <ArrowLeft className="h-4 w-4" />
                Back to Car Details
              </Button>
              <h1 className="text-2xl font-bold text-gray-900">
                Book Your Car
              </h1>
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Car Summary Card */}
            <div className="lg:col-span-1">
              <Card className="border-0 shadow-lg sticky top-8">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <CarIcon className="h-5 w-5 text-blue-600" />
                    Your Selection
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {/* Car Image */}
                  <div className="aspect-video relative overflow-hidden rounded-lg bg-gray-100">
                    {car.images && car.images.length > 0 ? (
                      <Image
                        src={car.images[0]}
                        alt={`${car.brand} ${car.model}`}
                        fill
                        className="object-cover"
                        onError={() => {}}
                      />
                    ) : (
                      <div className="flex items-center justify-center h-full">
                        <CarIcon className="h-12 w-12 text-gray-400" />
                      </div>
                    )}
                  </div>

                  {/* Car Details */}
                  <div>
                    <h3 className="font-bold text-lg text-gray-900">
                      {car.name || `${car.brand} ${car.model}`}
                    </h3>
                    <div className="flex items-center text-gray-600 text-sm mt-1">
                      <MapPin className="h-4 w-4 mr-1" />
                      {car.year} • {car.location_city}, {car.location_state}
                    </div>
                  </div>

                  {/* Pricing */}
                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <span className="text-gray-600">Daily Rate:</span>
                      <span className="font-semibold">
                        {formatPrice(car.price.rental_price_daily)}
                      </span>
                    </div>
                    {car.price.sale_price && car.price.sale_price > 0 && (
                      <div className="flex justify-between">
                        <span className="text-gray-600">Sale Price:</span>
                        <span className="font-semibold">
                          {formatPrice(car.price.sale_price)}
                        </span>
                      </div>
                    )}
                  </div>

                  <Badge variant="default" className="w-full justify-center">
                    Available
                  </Badge>
                </CardContent>
              </Card>
            </div>

            {/* Booking Form */}
            <div className="lg:col-span-2 space-y-6">
              {/* Customer Information */}
              <Card className="border-0 shadow-lg">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <User className="h-5 w-5 text-blue-600" />
                    Customer Information
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <Label>Name</Label>
                      <Input value={user?.username || ""} disabled />
                    </div>
                    <div>
                      <Label>Email</Label>
                      <Input value={user?.email || ""} disabled />
                    </div>
                  </div>
                  <div className="text-sm text-gray-600 bg-blue-50 p-3 rounded-lg">
                    <Info className="h-4 w-4 inline mr-1" />
                    Customer information is automatically filled from your
                    profile.
                  </div>
                </CardContent>
              </Card>

              {/* Booking Details */}
              <Card className="border-0 shadow-lg">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Calendar className="h-5 w-5 text-blue-600" />
                    Booking Details
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-6">
                  {/* Booking Type */}
                  <div>
                    <Label>Booking Type</Label>
                    <Select value={bookingType} onValueChange={setBookingType}>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="rental">Rental</SelectItem>
                        {car.price.sale_price && car.price.sale_price > 0 && (
                          <SelectItem value="purchase">Purchase</SelectItem>
                        )}
                      </SelectContent>
                    </Select>
                  </div>

                  {/* Date Selection for Rentals */}
                  {bookingType === "rental" && (
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <Label htmlFor="start-date">Start Date *</Label>
                        <Input
                          id="start-date"
                          type="date"
                          value={startDate}
                          onChange={(e) => setStartDate(e.target.value)}
                          min={new Date().toISOString().split("T")[0]}
                          required
                        />
                      </div>
                      <div>
                        <Label htmlFor="end-date">End Date *</Label>
                        <Input
                          id="end-date"
                          type="date"
                          value={endDate}
                          onChange={(e) => setEndDate(e.target.value)}
                          min={
                            startDate || new Date().toISOString().split("T")[0]
                          }
                          required
                        />
                      </div>
                    </div>
                  )}

                  {/* Notes */}
                  <div>
                    <Label htmlFor="notes">Additional Notes</Label>
                    <Textarea
                      id="notes"
                      placeholder="Any special requirements or notes..."
                      value={notes}
                      onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) =>
                        setNotes(e.target.value)
                      }
                      rows={3}
                    />
                  </div>
                </CardContent>
              </Card>

              {/* Cost Summary */}
              {((bookingType === "rental" && startDate && endDate) ||
                bookingType === "purchase") && (
                <Card className="border-0 shadow-lg">
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Calculator className="h-5 w-5 text-green-600" />
                      Cost Summary
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-3">
                      {bookingType === "rental" && (
                        <>
                          <div className="flex justify-between">
                            <span className="text-gray-600">
                              Daily Rate ({getDuration()}{" "}
                              {getDuration() === 1 ? "day" : "days"}):
                            </span>
                            <span className="font-medium">
                              {formatPrice(car.price.rental_price_daily)} ×{" "}
                              {getDuration()}
                            </span>
                          </div>
                          <div className="border-t pt-3">
                            <div className="flex justify-between items-center">
                              <span className="text-lg font-semibold">
                                Total Cost:
                              </span>
                              <span className="text-2xl font-bold text-green-600">
                                {formatPrice(calculateTotalCost())}
                              </span>
                            </div>
                          </div>
                        </>
                      )}

                      {bookingType === "purchase" && (
                        <div className="border-t pt-3">
                          <div className="flex justify-between items-center">
                            <span className="text-lg font-semibold">
                              Purchase Price:
                            </span>
                            <span className="text-2xl font-bold text-green-600">
                              {formatPrice(
                                car.price.sale_price ||
                                  car.price.rental_price_daily
                              )}
                            </span>
                          </div>
                        </div>
                      )}

                      <div className="text-sm text-gray-600 bg-yellow-50 p-3 rounded-lg">
                        <Info className="h-4 w-4 inline mr-1" />
                        Final pricing and payment will be handled by the car
                        owner.
                      </div>
                    </div>
                  </CardContent>
                </Card>
              )}

              {/* Submit Button */}
              <Card className="border-0 shadow-lg">
                <CardContent className="pt-6">
                  <Button
                    onClick={handleBooking}
                    disabled={
                      (bookingType === "rental" && (!startDate || !endDate)) ||
                      bookingLoading
                    }
                    size="lg"
                    className="w-full bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white font-semibold py-4 text-lg shadow-lg hover:shadow-xl transition-all duration-200"
                  >
                    {bookingLoading ? (
                      <>
                        <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-2" />
                        Processing...
                      </>
                    ) : (
                      <>
                        <CreditCard className="h-5 w-5 mr-2" />
                        {bookingType === "rental"
                          ? "Book Rental"
                          : "Submit Purchase Request"}
                      </>
                    )}
                  </Button>

                  <div className="text-center mt-4">
                    <p className="text-sm text-gray-600">
                      By proceeding, you agree to our terms of service and the
                      car owner will be notified of your booking request.
                    </p>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </div>
    </ProtectedRoute>
  );
}
