"use client";

import { motion } from "framer-motion";

const About = () => {
  return (
    <section id="about" className="py-20 bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
          className="text-center"
        >
          <h2 className="text-4xl md:text-5xl font-bold text-gray-900 mb-8">
            About{" "}
            <span className="bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">
              CarZone
            </span>
          </h2>
          <p className="text-xl text-gray-600 max-w-3xl mx-auto">
            With over 10 years of experience in the car rental industry, CarZone
            has established itself as the premier choice for luxury vehicle
            rentals. Our commitment to excellence and customer satisfaction sets
            us apart.
          </p>
        </motion.div>
      </div>
    </section>
  );
};

export default About;
