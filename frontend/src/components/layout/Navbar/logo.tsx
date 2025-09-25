import Image from "next/image"

export default function Logo() {
  return (
    <Image
      src="/logocar.png"
      alt="CarZone Logo"
      width={150}
      height={50}
      className="object-contain"
    />
  )
}
