import { useEffect, useState } from 'react'
import {
  ArrowDown,
  ArrowUp,
  BookOpen,
  Check,
  ChevronLeft,
  ChevronRight,
  Copy,
  Download,
  Plus,
  Search,
  Share,
  Smartphone,
  Star,
  X,
  Zap,
} from 'lucide-react'
import type { ReactNode } from 'react'

const NEVER_KEY = 'demso_install_never'

type InstallEvent = Event & {
  prompt: () => Promise<void>
  userChoice: Promise<{ outcome: 'accepted' | 'dismissed' }>
}

type Mode = 'ios' | 'inapp' | 'android'

function readNever(): boolean {
  try {
    return localStorage.getItem(NEVER_KEY) === '1'
  } catch {
    return false
  }
}

function writeNever() {
  try {
    localStorage.setItem(NEVER_KEY, '1')
  } catch {
    /* ignore */
  }
}

function isStandalone(): boolean {
  const nav = window.navigator as Navigator & { standalone?: boolean }
  return window.matchMedia('(display-mode: standalone)').matches || nav.standalone === true
}

function detectMode(): Mode {
  const ua = navigator.userAgent
  // Navigateurs intégrés (TikTok, Facebook, Instagram, WhatsApp...) : impossible d'installer depuis là
  if (/FBAN|FBAV|Instagram|TikTok|musical_ly|Bytedance|WhatsApp|Line\//i.test(ua)) return 'inapp'
  const iPadOS = navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1
  if (/iPhone|iPad|iPod/i.test(ua) || iPadOS) return 'ios'
  return 'android'
}

function Step({ n, icon, title, hint }: { n: number; icon: ReactNode; title: string; hint: string }) {
  return (
    <li className="flex items-center gap-4 rounded-2xl bg-gray-50 p-4">
      <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-gray-900 text-sm font-extrabold text-white">
        {n}
      </span>
      <span className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-[#F5B82E] text-xl font-extrabold text-gray-900">
        {icon}
      </span>
      <span className="min-w-0">
        <span className="block text-base font-extrabold leading-snug text-gray-900">{title}</span>
        <span className="block text-sm text-gray-600">{hint}</span>
      </span>
    </li>
  )
}

function Ring({ children }: { children: ReactNode }) {
  return (
    <span className="relative inline-flex">
      <span className="absolute -inset-2 animate-pulse rounded-full border-4 border-[#F5B82E]" />
      {children}
    </span>
  )
}

function GuideStep({ n, title, children }: { n: number; title: string; children: ReactNode }) {
  return (
    <li className="rounded-2xl bg-gray-50 p-4">
      <div className="flex items-center gap-3">
        <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-gray-900 text-sm font-extrabold text-white">
          {n}
        </span>
        <span className="text-base font-extrabold leading-snug text-gray-900">{title}</span>
      </div>
      <div className="mt-3">{children}</div>
    </li>
  )
}

function MenuRow({ icon, label, highlight }: { icon: ReactNode; label: string; highlight?: boolean }) {
  return (
    <div
      className={
        'flex items-center justify-between px-4 py-3 text-sm ' +
        (highlight ? 'border-2 border-[#F5B82E] bg-[#fdf1cc] font-extrabold text-gray-900' : 'text-gray-500')
      }
    >
      <span>{label}</span>
      {icon}
    </div>
  )
}

function IosGuide() {
  return (
    <ol className="mt-5 space-y-3">
      <GuideStep n={1} title="Appuie sur « Partager »">
        <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white">
          <div className="px-4 py-3 text-center text-xs text-gray-400">demso.com</div>
          <div className="flex items-center justify-around bg-gray-100 px-4 py-3 text-gray-400">
            <ChevronLeft size={22} />
            <ChevronRight size={22} />
            <Ring>
              <span className="flex h-9 w-9 items-center justify-center rounded-full bg-white text-blue-500">
                <Share size={22} />
              </span>
            </Ring>
            <BookOpen size={22} />
            <Copy size={22} />
          </div>
        </div>
        <p className="mt-3 flex items-center justify-center gap-2 text-sm font-extrabold text-gray-900">
          <ArrowUp size={18} /> Appuie ici, en bas de l'écran
        </p>
        <p className="mt-1 text-center text-xs text-gray-500">Tu ne vois pas ce bouton ? Appuie d'abord sur « ··· » en bas à droite.</p>
      </GuideStep>

      <GuideStep n={2} title="Descends et choisis « Sur l'écran d'accueil »">
        <p className="mb-2 flex items-center justify-center gap-2 text-sm font-extrabold text-gray-900">
          <ArrowDown size={18} /> Fais défiler le menu vers le bas
        </p>
        <div className="divide-y divide-gray-200 overflow-hidden rounded-2xl border border-gray-200 bg-white">
          <MenuRow label="Copier" icon={<Copy size={18} />} />
          <MenuRow label="Ajouter aux favoris" icon={<Star size={18} />} />
          <MenuRow
            highlight
            label="Sur l'écran d'accueil"
            icon={
              <span className="flex h-7 w-7 items-center justify-center rounded-md border-2 border-gray-900">
                <Plus size={16} />
              </span>
            }
          />
          <MenuRow label="Rechercher dans la page" icon={<Search size={18} />} />
        </div>
      </GuideStep>

      <GuideStep n={3} title="Appuie sur « Ajouter »">
        <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white">
          <div className="flex items-center justify-between bg-gray-100 px-4 py-3 text-sm">
            <span className="text-gray-400">Annuler</span>
            <Ring>
              <span className="rounded-full bg-white px-3 py-1 font-extrabold text-blue-600">Ajouter</span>
            </Ring>
          </div>
          <div className="flex items-center gap-3 px-4 py-4">
            <AppIcon />
            <div>
              <div className="font-extrabold text-gray-900">Demso</div>
              <div className="text-xs text-gray-500">demso.com</div>
            </div>
          </div>
        </div>
        <p className="mt-3 flex items-center justify-center gap-2 text-sm font-extrabold text-gray-900">
          <Check size={18} /> Demso est sur ton écran !
        </p>
      </GuideStep>
    </ol>
  )
}

function Steps({ mode }: { mode: Mode }) {
  if (mode === 'ios') return <IosGuide />
  if (mode === 'inapp') {
    return (
      <ol className="mt-5 space-y-3">
        <Step n={1} icon={<span>⋮</span>} title="Appuie sur les 3 points" hint="En haut à droite (⋮ ou ···)" />
        <Step n={2} icon={<Smartphone size={24} />} title="« Ouvrir dans le navigateur »" hint="Chrome ou Safari" />
        <Step n={3} icon={<Download size={24} />} title="Puis installe Demso" hint="Un bouton « Installer » apparaît" />
      </ol>
    )
  }
  return (
    <ol className="mt-5 space-y-3">
      <Step n={1} icon={<span>⋮</span>} title="Appuie sur les 3 points" hint="En haut à droite de Chrome" />
      <Step n={2} icon={<Download size={24} />} title="« Installer l'application »" hint="ou « Ajouter à l'écran d'accueil »" />
      <Step n={3} icon={<Check size={26} />} title="Confirme" hint="Demso est sur ton écran !" />
    </ol>
  )
}

function AppIcon() {
  return (
    <svg width="64" height="64" viewBox="0 0 100 100" aria-hidden="true" className="shrink-0 rounded-2xl">
      <rect width="100" height="100" fill="#faf6ea" />
      <g transform="translate(50 50) scale(1.7) translate(-24 -23)">
        <path d="M24 4 4 22h40z" fill="#F5B82E" />
        <rect x="9" y="22" width="30" height="20" rx="3" fill="#0E7C4A" />
        <rect x="20" y="29" width="8" height="13" rx="2" fill="#0B2E1F" />
      </g>
    </svg>
  )
}

export default function InstallPrompt() {
  const [installed, setInstalled] = useState<boolean>(() => isStandalone())
  const [evt, setEvt] = useState<InstallEvent | null>(null)
  const [open, setOpen] = useState(false)
  const [mode] = useState<Mode>(() => detectMode())

  useEffect(() => {
    if (installed) return

    const onBip = (e: Event) => {
      e.preventDefault()
      setEvt(e as InstallEvent)
    }
    const onInstalled = () => {
      setInstalled(true)
      setOpen(false)
      setEvt(null)
    }
    window.addEventListener('beforeinstallprompt', onBip)
    window.addEventListener('appinstalled', onInstalled)

    // À chaque entrée dans l'app, on propose (sauf si « Ne plus me proposer »)
    let timer: number | undefined
    if (!readNever()) timer = window.setTimeout(() => setOpen(true), 1500)

    return () => {
      window.removeEventListener('beforeinstallprompt', onBip)
      window.removeEventListener('appinstalled', onInstalled)
      if (timer) window.clearTimeout(timer)
    }
  }, [installed])

  if (installed) return null

  async function install() {
    if (!evt) return
    await evt.prompt()
    const choice = await evt.userChoice
    setEvt(null)
    if (choice.outcome === 'accepted') setOpen(false)
  }

  function never() {
    writeNever()
    setOpen(false)
  }

  return (
    <>
      {/* Bouton toujours visible, même après « Ne plus me proposer » */}
      {!open && (
        <button
          type="button"
          onClick={() => setOpen(true)}
          className="fixed bottom-4 right-4 z-40 inline-flex items-center gap-2 rounded-full bg-gray-900 px-4 py-3 text-sm font-semibold text-white shadow-lg hover:bg-black"
        >
          <Download size={18} />
          Installer l'app
        </button>
      )}

      {open && (
        <div className="fixed inset-0 z-50 flex items-end justify-center bg-black/50" onClick={() => setOpen(false)}>
          <div
            role="dialog"
            aria-modal="true"
            aria-label="Installer Demso"
            onClick={(e) => e.stopPropagation()}
            className="max-h-[92vh] w-full max-w-md overflow-y-auto rounded-t-3xl bg-white p-6 pb-[max(1.5rem,env(safe-area-inset-bottom))] shadow-2xl"
          >
            <div className="flex items-start justify-between gap-3">
              <div className="flex items-center gap-4">
                <AppIcon />
                <div>
                  <h2 className="text-xl font-extrabold leading-tight text-gray-900">Installe Demso sur ton téléphone</h2>
                  <p className="mt-1 text-sm text-gray-600">Gratuit. Ouvre Demso en un seul toucher.</p>
                </div>
              </div>
              <button
                type="button"
                onClick={() => setOpen(false)}
                aria-label="Fermer"
                className="rounded-full p-2 text-gray-500 hover:bg-gray-100"
              >
                <X size={20} />
              </button>
            </div>

            <div className="mt-5 grid grid-cols-3 gap-2 text-center text-xs font-semibold text-gray-800">
              <div className="rounded-2xl bg-brand-50 p-3">
                <Smartphone className="mx-auto mb-1 text-brand-700" size={22} />
                Icône sur ton écran
              </div>
              <div className="rounded-2xl bg-brand-50 p-3">
                <Zap className="mx-auto mb-1 text-brand-700" size={22} />
                Ouvre en 1 toucher
              </div>
              <div className="rounded-2xl bg-brand-50 p-3">
                <Search className="mx-auto mb-1 text-brand-700" size={22} />
                Plus besoin de chercher
              </div>
            </div>

            {evt ? (
              <button type="button" onClick={install} className="btn btn-primary mt-6 w-full">
                <Download size={20} />
                Installer
              </button>
            ) : (
              <Steps mode={mode} />
            )}

            <button
              type="button"
              onClick={() => setOpen(false)}
              className="btn btn-outline mt-4 w-full"
            >
              Plus tard
            </button>
            <button type="button" onClick={never} className="mt-3 w-full text-center text-sm text-gray-500 underline">
              Ne plus me proposer
            </button>
          </div>
        </div>
      )}
    </>
  )
}
