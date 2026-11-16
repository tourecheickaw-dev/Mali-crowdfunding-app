import { Link } from 'react-router-dom'
import { UserRound } from 'lucide-react'
import Logo from './Logo'

export default function Header() {
  return (
    <header className="sticky top-0 z-10 border-b border-gray-200 bg-white">
      <div className="mx-auto flex h-16 max-w-5xl items-center justify-between px-4">
        <Link to="/" aria-label="Demso, accueil">
          <Logo size={34} />
        </Link>

        <nav className="flex items-center gap-2 sm:gap-3">
          <Link
            to="/cagnottes"
            className="hidden px-3 py-2 font-medium text-gray-700 transition hover:text-[#111827] sm:block"
          >
            Explorer
          </Link>

         <Link
  to="/mon-espace"
  aria-label="Mon espace"
  className="inline-flex h-11 items-center justify-center gap-2 rounded-full bg-brand-50 px-3 text-sm font-semibold text-brand-800 hover:bg-brand-100 sm:px-4"
>
  <UserRound size={20} />
  <span className="hidden sm:inline">Mon espace</span>
</Link>

          <Link
            to="/creer"
        className="rounded-full bg-gray-900 px-5 py-2.5 text-sm font-semibold text-white hover:bg-black"
          >
            Créer une cagnotte
          </Link>
        </nav>
      </div>
    </header>
  )
}