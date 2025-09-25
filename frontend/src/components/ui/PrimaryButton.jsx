/* eslint-disable react/prop-types */
// import React from 'react'

const Button = ({ title, id, leftIcon, containerClass, rightIcon, onClick }) => {
    return (
        <button id={id} className={` group relative z-10 w-fit cursor-pointer overflow-hidden rounded-full px-7 py-3 bg-primary text-foreground ${containerClass}`} onClick={onClick}>
            {leftIcon}
            <span className="relative inline-flex overflow-hidden font-general text-xs uppercase ">
                <div className="translate-y-0 skew-y-0 transition duration-500 group-hover:translate-y-[-160%] group-hover:skew-y-12">
                    {title}
                </div>
                <div className="absolute translate-y-[164%] skew-y-12 transition duration-500 group-hover:translate-y-0 group-hover:skew-y-0">
                    {title}
                </div>
            </span>
            {rightIcon}
        </button>
    )
}

export default Button