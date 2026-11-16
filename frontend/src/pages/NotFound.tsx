import { Link } from 'react-router-dom'

export default function NotFound() {
  return (
    <div className="py-24 text-brand-600enter">
      <p className="text-6xl font-extrabold text-brand-600rand-900rand-600">404</p>
      <p className="mt-3 text-xl text-gray-700">Cette page n'existe pas.</p>
      <Link to="/" className="mt-6 inline-block rounded-full bg-[#111827] px-8 py-4 text-lg font-semibold text-white">
        Retour à l'accueil
      </Link>
    </div>
  )
}