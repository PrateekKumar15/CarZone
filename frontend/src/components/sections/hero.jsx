// import React from 'react'
import { useState, useRef, useEffect } from "react"
import Button from "@/components/ui/PrimaryButton"
import { TiLocationArrow } from "react-icons/ti"
import { useGSAP } from "@gsap/react"
import gsap from "gsap"
import ScrollTrigger from "gsap/ScrollTrigger"
gsap.registerPlugin(ScrollTrigger)
const Hero = () => {
    // const [currentIndex, setCurrentIndex] = useState(1)
    // const [hasClicked, setHasClicked] = useState(false)
    const [isloading, setIsLoading] = useState(true);
    const [loadedVideos, setLoadedVideos] = useState(0);
    const totalVideos = 4;
    const nextVideoRef = useRef(null);
    // const upcomingVideoIndex = ((currentIndex % totalVideos + 1));
    // const getVideoSrc = (index) => {
    //     return `./videos/hero-${index}.mp4`
    // }

    const handleVideoLoaded = () => {
        setLoadedVideos((prev) => (prev + 1));
    }



    // const handleMiniVdClick = () => {
    //     setHasClicked(true);
    //     setCurrentIndex(upcomingVideoIndex);

    // }

    useEffect(() => {
        loadedVideos == totalVideos - 1;
        setIsLoading(false);
    }, [loadedVideos]);


    // useGSAP(() => {
    //     if (hasClicked) {
    //         gsap.set('#next-video', {
    //             visibility: 'visible'
    //         })

    //         gsap.to('#next-video', {
    //             transformOrigin: 'center center',
    //             scale: 1,
    //             width: '100%',
    //             height: '100%',
    //             duration: 1,
    //             ease: 'power1.inOut',
    //             onStart: () => {
    //                 nextVideoRef.current.play();
    //             }
    //         })

    //         gsap.from('#current-video', {
    //             transformOrigin: 'center center',
    //             scale: 0,
    //             duration: 1.5,
    //             ease: 'power1.inOut',
    //         })

    //     }

    // }
    //     , { dependencies: [currentIndex], revertOnUpdate: true })


    useGSAP(() => {
        gsap.set('#video-frame', {
            clipPath: "polygon(14% 0, 72% 0, 90% 90%, 0 100%)",
            borderRadius: '0 0 40% 10%'
        })

        gsap.from("#video-frame", {
            clipPath: "polygon(0% 0%, 100% 0%, 100% 100%, 0% 100%)",
            borderRadius: "0% 0% 0% 0%",
            ease: "power1.inOut",
            scrollTrigger: {
                trigger: "#video-frame",
                start: "center center",
                end: "bottom center",
                scrub: true,
            },
        });
    })


    return (
        <div className="relative h-dvh w-screen overflow-x-hidden">
            {isloading && (
                <div className="flex-center absolute h-dvh w-screen z-[100]  overflow-hidden bg-violet-50">
                    <div className="three-body">
                        <div className="three-body__dot"></div>
                        <div className="three-body__dot"></div>
                        <div className="three-body__dot"></div>
                    </div>
                </div>
            )}
            <div id="video-frame" className="relative z-10 h-dvh w-screen overflow-hidden rounded-lg" >
                <div>
                    {/* <div className="mask-clip-path absolute-center z-50 size-64 cursor-pointer overflow-hidden rounded-lg">
                        <div onClick={handleMiniVdClick} className="origin-center scale-50 opacity-0 transition-all duration-500 hover:scale-100 hover:opacity-100" >
                            <video ref={nextVideoRef} src={getVideoSrc(upcomingVideoIndex)} loop muted id="current-video" className="size-64 origin-center scale-150 object-cover object-center" onLoadedData={handleVideoLoaded}></video>
                        </div>

                    </div> */}

                    {/* <video ref={nextVideoRef} src="./videos/hero-4.mp4" loop muted id="next-video" className="absolute-center invisible z-20 size-64 object-cover object-center"
                        onLoadedData={handleVideoLoaded}
                    ></video> */}
                    <video src="./videos/hero-4.mp4"
                        autoPlay loop muted className="absolute left-0 top-0 size-full object-cover object-center" onLoadedData={handleVideoLoaded}
                    ></video>

                </div>
                <h1 className="special-font hero-heading absolute bottom-5 right-5 z-40 text-primary">
                    R<b>e</b>nt or B<b>u</b>y
                </h1>
                <div className="absolute left-0 top-auto z-40 size-full">
                    <div className="mt-24 px-5 sm:px-10">
                        <h1 className="special-font hero-heading block text-accent-foreground ">c<b>a</b>rzo<b>n</b>e</h1>
                        <p className="mb-5 max-w-64 font-robert-regular text-accent-foreground ">Enter the MetaGame Layer<br /> Unleash Economy</p>
                        <Button id='watch-trailer' title='View All Vehicles' onClick={()=>{
                            window.location.href = '/cars';
                        }} leftIcon={<TiLocationArrow />} containerClass="bg-primary text-primary-foreground flex-center gap-1 " />

                    </div>
                </div>
            </div>

            <h1 className="special-font hero-heading absolute bottom-5 right-5  text-foreground">
                R<b>e</b>nt or B<b>u</b>y
            </h1>
        </div>
    )
}

export default Hero