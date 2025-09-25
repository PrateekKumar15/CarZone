"use client";

import { useEffect } from "react";
import Hero from "@/components/sections/hero";
import FeaturedCars from "@/components/sections/featured-cars";
import Services from "@/components/sections/services";
import About from "@/components/sections/about";
import Testimonials from "@/components/sections/testimonials";
import Contact from "@/components/sections/contact";
import Footer from "@/components/layout/footer";
import Navbar from "@/components/layout/Navbar/Navbar";

export default function HomePage() {
  useEffect(() => {
    // Initialize any global animations or effects here
    document.body.style.overflowX = "hidden";

    return () => {
      document.body.style.overflowX = "auto";
    };
  }, []);

  return (
    <main className="min-h-screen bg-background">
      <div className="fixed top-0 left-0 w-full z-50 bg-background backdrop-blur-md border-b">
        <Navbar />
      </div>

      <Hero />
      <FeaturedCars />
      <Services />
      <About />
      <Testimonials />
      <Contact />
      <Footer />
    </main>
  );
}
