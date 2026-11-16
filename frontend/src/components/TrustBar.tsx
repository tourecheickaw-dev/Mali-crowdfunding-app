import { Smartphone, Receipt, MessagesSquare } from 'lucide-react'

const items = [
  {
    icon: Smartphone,
    title: 'Paiement simple',
    text: 'Wave, Orange Money, Moov ou carte bancaire.',
  },
  {
    icon: Receipt,
    title: 'Frais clairs',
    text: 'Le montant des frais est affiché avant de payer.',
  },
  {
    icon: MessagesSquare,
    title: 'Messages visibles',
    text: "Chaque don peut avoir son message, et l'organisateur peut répondre.",
  },
]

export default function TrustBar() {
  return (
    <section className="mx-auto max-w-6xl px-4 pb-8 pt-28">
      <div className="grid gap-6 md:grid-cols-3">
        {items.map(({ icon: Icon, title, text }) => (
          <div
            key={title}
            className="rounded-3xl bg-white p-8 shadow-[0_10px_40px_-15px_rgba(0,0,0,0.2)] ring-1 ring-gray-200 transition duration-200 hover:-translate-y-1 hover:shadow-[0_20px_50px_-15px_rgba(0,0,0,0.25)]"
          >
           <span className="flex h-14 w-14 items-center justify-center rounded-2xl bg-[#F5B82E] text-gray-900">
              <Icon size={26} strokeWidth={2.2} />
            </span>
            <h3 className="mt-6 text-xl font-extrabold tracking-tight text-gray-900">{title}</h3>
            <p className="mt-2 leading-relaxed text-gray-600">{text}</p>
          </div>
        ))}
      </div>
    </section>
  )
}