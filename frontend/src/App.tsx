import { Routes, Route } from 'react-router-dom'
import Header from './components/Header'
import Footer from './components/Footer'
import Home from './pages/Home'
import Explore from './pages/Explore'
import CampaignPage from './pages/CampaignPage'
import NotFound from './pages/NotFound'
import DonatePage from './pages/DonatePage'
import CreateCampaignPage from './pages/CreateCampaignPage'
import MyCampaignsPage from './pages/MyCampaignsPage'
import OrganizerLayout from './pages/OraganizerLayout'
import MessagesPage from './pages/MessagesPage'
import PayoutPage from './pages/PayoutPage'
import { About, Terms, Privacy, LegalNotice, Contact } from './pages/InfoPages'
import InstallPrompt from './components/InstallPrompt'

export default function App() {
  return (
    <>
      <Header />
      <main className="mx-auto max-w-5xl px-4">
        <Routes>
                    <Route path="/mon-espace" element={<MyCampaignsPage />} />
          <Route path="/" element={<Home />} />
          <Route path="/cagnottes" element={<Explore />} />
          <Route path="/c/:slug" element={<CampaignPage />} />
          <Route path="*" element={<NotFound />} />
                    <Route path="/c/:slug/don" element={<DonatePage />} />
          <Route path="/creer" element={<CreateCampaignPage />} />
          <Route path="/mon-espace" element={<MyCampaignsPage />} />
  <Route path="/mon-espace/layout" element={<OrganizerLayout />} />
  <Route path="/mon-espace/c/:id/messages" element={<MessagesPage />} />
  <Route path="/mon-espace/payouts" element={<PayoutPage />} />
  <Route path="/a-propos" element={<About />} />
<Route path="/conditions" element={<Terms />} />
<Route path="/confidentialite" element={<Privacy />} />
<Route path="/mentions-legales" element={<LegalNotice />} />
<Route path="/contact" element={<Contact />} />

  <Route path="*" element={<NotFound />} />
</Routes>
        
      </main>
      <Footer />
      <InstallPrompt />
    </>
  )
}