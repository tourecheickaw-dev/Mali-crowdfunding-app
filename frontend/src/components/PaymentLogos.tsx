const logos = [
  { name: 'Wave', src: '/logos/wave.png' },
  { name: 'Orange Money', src: '/logos/orange-money.png' },
  { name: 'Moov Money', src: '/logos/moov-money.png' },
]

const tile =
  'flex h-16 w-28 items-center justify-center rounded-2xl border border-gray-200 bg-white shadow-sm transition hover:-translate-y-0.5 hover:shadow-md'

export default function PaymentLogos() {
  return (
    <section className="mt-20 text-brand-600enter">
      

      <div className="mt-6 flex flex-wrap items-center justify-center gap-3 sm:gap-4">
        {logos.map((l) => (
          <div key={l.name} className={tile} title={l.name}>
            <img src={l.src} alt={l.name} className="h-9 w-auto max-w-20 rounded-md object-contain" />
          </div>
        ))}

        <div className={tile} title="Visa">
          <span className="text-brand-600rand-900xl font-extrabold italic tracking-tight text-[#1a1f71]">VISA</span>
        </div>

        <div className={tile} title="Mastercard">
          <svg width="46" height="30" viewBox="0 0 46 30" aria-label="Mastercard">
            <circle cx="16" cy="15" r="13" fill="#eb001b" />
            <circle cx="30" cy="15" r="13" fill="#f79e1b" />
            <path d="M23 4.2a13 13 0 0 1 0 21.6 13 13 0 0 1 0-21.6z" fill="#ff5f00" />
          </svg>
        </div>
      </div>
    </section>
  )
}