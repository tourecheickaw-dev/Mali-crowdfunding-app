export default function HowItWorks() {
  const steps = [
    {
      title: 'Crée ta cagnotte en quelques minutes',
      text: 'Ajoute un titre, une photo et ton objectif. Tu peux modifier ta cagnotte à tout moment.',
    },
    {
      title: 'Partage-la avec ton entourage',
      text: 'Envoie le lien sur WhatsApp, Facebook, TikTok ou Instagram. Les dons arrivent de partout.',
    },
    {
      title: "Reçois l'argent sur Wave ou Orange Money",
      text: 'Les dons sont versés sur ton compte mobile money. Tu suis tout en direct.',
    },
  ]

  return (
    <section className="mx-auto max-w-6xl px-4 py-16">
      <div className="grid items-center gap-12 md:grid-cols-2">
        <div className="mx-auto w-full max-w-sm rounded-[2rem] bg-gradient-to-br from-orange-400 to-orange-500 p-5 shadow-xl">
          <div className="rounded-2xl bg-white p-5 shadow-sm">
            <p className="text-xs font-bold uppercase tracking-wider text-gray-900">Exemple</p>
            <p className="mt-1 text-lg font-extrabold text-gray-900">Ma cagnotte</p>
            <p className="mt-5 text-2xl font-extrabold text-gray-900">125 000 FCFA</p>
            <p className="text-sm text-gray-900">sur 200 000 FCFA</p>
            <div className="mt-3 h-2 overflow-hidden rounded-full bg-gray-200">
              <div className="h-full w-[62%] rounded-full bg-brand-600" />
            </div>
            <div className="mt-5 flex items-center gap-3">
              <span className="flex h-9 w-9 items-center justify-center rounded-full bg-brand-100 text-sm font-bold text-brand-700">
                A
              </span>
              <div>
                <p className="text-sm font-semibold text-gray-900">Un donateur</p>
                <p className="text-xs text-gray-900">« Courage à toute la famille »</p>
              </div>
            </div>
          </div>
        </div>

        <div>
          <h2 className="text-3xl font-extrabold leading-tight tracking-tight text-gray-900">
            Collecter de l'argent avec Demso, c'est simple et sûr
          </h2>
          <ol className="mt-8 space-y-6">
            {steps.map((s, i) => (
              <li key={s.title} className="flex gap-4 border-b border-gray-200 pb-6 last:border-0">
                <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full border-2 border-gray-900 text-base font-bold text-gray-900">
                  {i + 1}
                </span>
                <div>
                  <p className="text-lg font-bold text-gray-900">{s.title}</p>
                  <p className="mt-1 text-gray-900">{s.text}</p>
                </div>
              </li>
            ))}
          </ol>
        </div>
      </div>
    </section>
  )
}