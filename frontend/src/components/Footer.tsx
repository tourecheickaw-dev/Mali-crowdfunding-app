import { Link } from 'react-router-dom'
import Logo from './Logo'

const links = [
  { to: '/a-propos', label: 'À propos' },
  { to: '/contact', label: 'Contact' },
  { to: '/conditions', label: "Conditions d'utilisation" },
  { to: '/confidentialite', label: 'Confidentialité' },
  { to: '/mentions-legales', label: 'Mentions légales' },
]

export default function Footer() {
  return (
    <footer className="mt-20 bg-gray-100">
      <div className="mx-auto max-w-5xl px-4 pb-8 pt-12">
        <Logo size={30} />

        <nav className="mt-6 flex flex-wrap gap-x-8 gap-y-3">
          {links.map((l) => (
            <Link
              key={l.to}
              to={l.to}
              className="font-medium text-gray-700 transition hover:text-gray-900 hover:underline"
            >
              {l.label}
            </Link>
          ))}
        </nav>

        <p className="mt-10 border-t border-gray-300 pt-6 text-sm text-gray-600">
          © 2026 Demso. Tous droits réservés.
        </p>
      </div>
    </footer>
  )
}