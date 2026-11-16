type Props = { size?: number; className?: string }

export default function Logo({ size = 28, className = '' }: Props) {
  return (
    <span
      className={`inline-block whitespace-nowrap font-extrabold leading-none tracking-tight text-brand-900 ${className}`}
      style={{ fontSize: size }}
    >
      <svg viewBox="0 0 100 100" aria-hidden="true" className="mr-[0.18em] inline-block h-[1em] w-[1em] align-baseline">
        <path d="M50 2 C53 20 76 41 99 54 C70 61 30 61 1 54 C24 41 47 20 50 2Z" fill="#F5B82E" />
        <path d="M18 56 H82 V86 Q82 97 71 97 H29 Q18 97 18 86Z" fill="#0E7C4A" />
        <path d="M40 97 V74 a10 10 0 0 1 20 0 V97Z" fill="#0A5C37" />
      </svg>
      Demso
    </span>
  )
}