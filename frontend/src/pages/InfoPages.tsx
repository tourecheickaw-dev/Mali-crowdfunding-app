import { useEffect, type ReactNode } from 'react'
import { Link } from 'react-router-dom'

// À MODIFIER : ton vrai e-mail de contact (et ton numéro WhatsApp si tu veux en mettre un).
const CONTACT_EMAIL = 'demso.support@gmail.com'
const UPDATED = '3 novembre 2026'

function Shell({ title, intro, children }: { title: string; intro?: string; children: ReactNode }) {
  useEffect(() => {
    window.scrollTo(0, 0)
  }, [])
  return (
    <div className="mx-auto max-w-3xl py-10">
      <h1 className="text-3xl font-extrabold leading-tight tracking-tight text-gray-900 sm:text-4xl">{title}</h1>
      {intro && <p className="mt-3 text-lg text-gray-700">{intro}</p>}
      <div className="mt-8 space-y-8">{children}</div>
      <p className="mt-12 text-sm text-gray-500">Dernière mise à jour : {UPDATED}</p>
    </div>
  )
}

function Block({ title, children }: { title: string; children: ReactNode }) {
  return (
    <section className="rounded-3xl bg-white p-6 shadow-sm ring-1 ring-gray-200">
      <h2 className="text-xl font-extrabold text-gray-900">{title}</h2>
      <div className="mt-3 space-y-3 leading-relaxed text-gray-800">{children}</div>
    </section>
  )
}

export function About() {
  return (
    <Shell
      title="À propos de Demso"
      intro="Demso veut que chacun puisse aider ses proches, où qu'il soit, simplement et en confiance."
    >
      <Block title="Pourquoi Demso ?">
        <p>
          Aujourd'hui, pour aider un proche au pays, on partage un numéro de téléphone sur WhatsApp ou TikTok, et on
          espère que tout le monde fera confiance. Demso rend cela plus simple et plus clair : une page pour chaque
          cagnotte, un lien à partager, et des dons qui arrivent directement chez la personne qui organise.
        </p>
      </Block>
      <Block title="Le nom">
        <p>« Demso » vient du bambara : « dèmè » (aider) et « so » (la maison). C'est la maison de l'entraide.</p>
      </Block>
      <Block title="Comment ça marche">
        <p>
          Tu crées ta cagnotte en quelques minutes, tu la partages avec ton entourage, et les dons arrivent sur ton
          compte Wave, Orange Money, Moov ou bancaire. Les donateurs peuvent payer avec ces moyens ou par carte
          bancaire, depuis n'importe quel pays.
        </p>
      </Block>
      <Block title="Une question ?">
        <p>
          Écris-nous à <a className="font-semibold underline" href={`mailto:${CONTACT_EMAIL}`}>{CONTACT_EMAIL}</a>.
        </p>
      </Block>
    </Shell>
  )
}

export function Terms() {
  return (
    <Shell
      title="Conditions d'utilisation"
      intro="Merci de les lire avant de créer une cagnotte ou de faire un don."
    >
      <Block title="1. Ce qu'est Demso">
        <p>
          Demso est une plateforme qui permet à une personne (l'organisateur) de présenter une cagnotte et de recevoir
          des dons. Demso n'est ni une association, ni une banque. Demso ne garantit pas que l'organisateur utilisera
          les fonds comme annoncé.
        </p>
      </Block>
      <Block title="2. Les dons">
        <p>
          Les dons sont envoyés par un prestataire de paiement vers le compte indiqué par l'organisateur. Demso ne
          garde que ses frais de service, affichés avant le paiement (actuellement 6 %, ajoutés au montant du don).
        </p>
        <p>
          Un don est en principe définitif. En cas de fraude avérée ou d'erreur manifeste, contactez-nous : nous
          étudierons chaque situation.
        </p>
      </Block>
      <Block title="3. Obligations de l'organisateur">
        <ul className="list-disc space-y-2 pl-6">
          <li>Donner des informations vraies et ne pas tromper les donateurs.</li>
          <li>Utiliser l'argent reçu pour le but annoncé.</li>
          <li>Avoir le droit d'utiliser les photos et textes publiés, et l'accord des personnes montrées.</li>
          <li>Respecter la dignité des personnes, en particulier sur les photos médicales ou difficiles.</li>
          <li>Indiquer un compte de réception qui lui appartient ou dont il a l'autorisation.</li>
        </ul>
      </Block>
      <Block title="4. Contenus interdits">
        <p>
          Sont interdits : les fausses cagnottes, les arnaques, les contenus violents, haineux ou illégaux, et toute
          collecte pour une activité interdite par la loi.
        </p>
      </Block>
      <Block title="5. Modération">
        <p>
          Demso peut masquer ou supprimer une cagnotte, un message ou un compte qui ne respecte pas ces conditions,
          ou en cas de doute sérieux. Un message peut aussi être masqué par l'organisateur de la cagnotte.
        </p>
      </Block>
      <Block title="6. Responsabilité">
        <p>
          Demso met tout en œuvre pour que le service fonctionne, mais ne peut pas garantir qu'il soit toujours
          disponible. Les opérateurs de paiement (Wave, Orange Money, Moov, banques, cartes) peuvent avoir des
          interruptions ou des plafonds, indépendants de Demso.
        </p>
      </Block>
      <Block title="7. Évolutions et contact">
        <p>
          Ces conditions peuvent évoluer ; la date de mise à jour est indiquée en bas de page. Pour toute question :{' '}
          <a className="font-semibold underline" href={`mailto:${CONTACT_EMAIL}`}>{CONTACT_EMAIL}</a>.
        </p>
      </Block>
    </Shell>
  )
}

export function Privacy() {
  return (
    <Shell
      title="Confidentialité et données personnelles"
      intro="Voici, simplement, quelles informations Demso utilise et pourquoi."
    >
      <Block title="Les informations que nous collectons">
        <ul className="list-disc space-y-2 pl-6">
          <li>Ton nom et ton numéro de téléphone (vérifié par Wave ou Orange Money).</li>
          <li>Ton e-mail, seulement si tu choisis de le donner.</li>
          <li>Les informations de ta cagnotte : titre, texte, photos, objectif.</li>
          <li>Pour les organisateurs : les comptes de réception (nom, numéro de téléphone ou de compte).</li>
          <li>Pour les donateurs : le montant du don, ton éventuel message, et ton nom si tu ne restes pas anonyme.</li>
          <li>Des données techniques (par exemple l'adresse IP et les journaux de connexion) pour la sécurité.</li>
        </ul>
      </Block>
      <Block title="Ce que nous ne voyons jamais">
        <p>
          Demso ne voit ni ne conserve les numéros de carte bancaire, les codes secrets, les codes reçus par SMS ni
          les mots de passe Wave ou Orange Money. Ils sont traités directement par le prestataire de paiement, la
          banque ou l'opérateur.
        </p>
      </Block>
      <Block title="À quoi servent ces informations">
        <ul className="list-disc space-y-2 pl-6">
          <li>Faire fonctionner les cagnottes et envoyer les dons au bon compte.</li>
          <li>Afficher publiquement la cagnotte, les dons et les messages (sans ton numéro ni ton e-mail).</li>
          <li>Protéger la plateforme contre la fraude et les abus.</li>
          <li>T'écrire si c'est nécessaire au service.</li>
        </ul>
        <p>Nous ne vendons pas tes données.</p>
      </Block>
      <Block title="Ce qui est public">
        <p>
          Le titre, la description, les photos, le nom de l'organisateur, les montants et les messages d'une cagnotte
          sont visibles par tous. Ton numéro de téléphone, ton e-mail et les comptes de réception ne sont jamais
          affichés publiquement. Un don anonyme n'affiche pas ton nom.
        </p>
      </Block>
      <Block title="Hébergement et partage">
        <p>
          Les données sont hébergées chez Amazon Web Services, en Europe (Paris). Elles sont partagées uniquement avec
          les prestataires nécessaires au service (paiement, hébergement) et si la loi nous y oblige.
        </p>
      </Block>
      <Block title="Tes droits">
        <p>
          Tu peux demander à voir, corriger ou supprimer tes données. Écris-nous à{' '}
          <a className="font-semibold underline" href={`mailto:${CONTACT_EMAIL}`}>{CONTACT_EMAIL}</a>. Certaines
          données (par exemple les traces des paiements) doivent être conservées plus longtemps pour des raisons
          légales ou de sécurité.
        </p>
      </Block>
      <Block title="Cookies">
        <p>Demso n'utilise pas de publicité ciblée. Seules des données techniques nécessaires au fonctionnement sont enregistrées dans ton navigateur.</p>
      </Block>
    </Shell>
  )
}

export function LegalNotice() {
  return (
    <Shell title="Mentions légales">
     <Block title="Éditeur du site">
        <p>Demso — entreprise individuelle</p>
        <p>Activité : édition de logiciels (application en ligne).</p>
        <p>Immatriculée au RCCM sous le numéro ML-BKO-01-2026-A-07849 (Bamako, Mali).</p>
        <p>Siège : Moribabougou, Commune I, Bamako, Mali.</p>
        <p>NIF : en cours d'attribution.</p>
        <p>
          Contact : <a className="font-semibold underline" href="mailto:demso.support@gmail.com">demso.support@gmail.com</a>
        </p>
      </Block>
      <Block title="Directeur de la publication">
        <p>Cheick Abdoul Wahab Touré</p>
      </Block>
      <Block title="Hébergement">
        <p>
          Amazon Web Services EMEA SARL — région Europe (Paris). Le site est diffusé par un hébergeur de sites web
          (à indiquer quand le choix sera fait).
        </p>
      </Block>
    </Shell>
  )
}

export function Contact() {
  return (
    <Shell title="Nous contacter" intro="Une question, un problème, une idée ? Écris-nous.">
      <Block title="Par e-mail">
        <p>
          <a className="text-lg font-semibold underline" href={`mailto:${CONTACT_EMAIL}`}>{CONTACT_EMAIL}</a>
        </p>
        <p>Nous répondons dès que possible.</p>
      </Block>
      <Block title="Avant de nous écrire">
        <p>
          Pour un problème avec un don, indique le nom de la cagnotte, le montant, la date et le moyen de paiement.
          Tu peux aussi regarder les{' '}
          <Link to="/conditions" className="font-semibold underline">conditions d'utilisation</Link>.
        </p>
      </Block>
    </Shell>
  )
}
