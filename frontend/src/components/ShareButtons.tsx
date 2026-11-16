import { useState } from 'react'
import { Check, Link2, MessageCircle, Share2 } from 'lucide-react'

export default function ShareButtons({ title }: { title: string }) {
  const [copied, setCopied] = useState(false)
  const url = window.location.href
  const text = `Aidez-nous : ${title}`
  const canNativeShare = typeof navigator.share === 'function'

  const whatsapp = `https://wa.me/?text=${encodeURIComponent(text + '\n' + url)}`
  const facebook = `https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(url)}`

  async function nativeShare() {
    try {
      await navigator.share({ title, text, url })
    } catch {
      // l'utilisateur a fermé le menu : rien à faire
    }
  }

  async function copyLink() {
    try {
      await navigator.clipboard.writeText(url)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch {
      window.prompt('Copiez ce lien :', url)
    }
  }

  const base =
    'flex h-12 items-center justify-center gap-2 rounded-2xl px-3 text-sm font-semibold transition hover:-translate-y-0.5 hover:shadow-md'

  return (
    <div className="grid grid-cols-2 gap-3">
      <a
        href={whatsapp}
        target="_blank"
        rel="noopener noreferrer"
        className={`${base} text-white`}
        style={{ backgroundColor: '#25D366' }}
      >
        <MessageCircle size={20} />
        WhatsApp
      </a>

      <a
        href={facebook}
        target="_blank"
        rel="noopener noreferrer"
        className={`${base} text-white`}
        style={{ backgroundColor: '#1877F2' }}
      >
        <span className="text-xl font-extrabold leading-none">f</span>
        Facebook
      </a>

      {canNativeShare && (
        <button type="button" onClick={nativeShare} className={`${base} bg-gray-900 text-white`}>
          <Share2 size={20} />
          TikTok, Instagram…
        </button>
      )}

      <button
        type="button"
        onClick={copyLink}
        className={`${base} border-2 ${
          copied ? 'border-brand-600 bg-brand-50 text-brand-700' : 'border-gray-900 bg-white text-gray-900'
        } ${canNativeShare ? '' : 'col-span-2'}`}
      >
        {copied ? <Check size={20} /> : <Link2 size={20} />}
        {copied ? 'Lien copié' : 'Copier le lien'}
      </button>
    </div>
  )
}
