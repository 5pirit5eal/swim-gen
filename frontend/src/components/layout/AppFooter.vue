<script setup lang="ts">
import BaseModal from '@/components/ui/BaseModal.vue'
import PrivacyPolicy from '../legal/PrivacyPolicy.vue'
import { ref } from 'vue'
// Footer component with privacy links and legal information
// Handlers for modal dialogs will be implemented when needed

const showPrivacyPolicyModal = ref(false)

const showTerms = () => {
  // TODO: Implement terms modal/page
  console.log('Terms clicked')
}

const showImpressum = () => {
  // TODO: Implement impressum modal/page (required for German/EU law)
  console.log('Impressum clicked')
}

const showContact = () => {
  // TODO: Implement contact modal/page
  console.log('Contact clicked')
}
</script>

<template>
  <footer class="app-footer">
    <div class="footer-container">
      <div class="footer-content">
        <div class="footer-section">
          <h3>Swim RAG</h3>
          <p>AI-powered training plan generator for swimmers and triathletes.</p>
        </div>

        <div class="footer-section">
          <h4>Privacy & Legal</h4>
          <ul class="footer-links">
            <li>
              <a href="#privacy" id="privacy-policy" @click="showPrivacyPolicyModal = true"
                >Privacy Policy</a
              >
            </li>
            <li><a href="#terms" @click.prevent="showTerms">Terms of Service</a></li>
            <li><a href="#impressum" @click.prevent="showImpressum">Impressum</a></li>
          </ul>
        </div>

        <div class="footer-section">
          <h4>About</h4>
          <ul class="footer-links">
            <li><a href="#contact" @click.prevent="showContact">Contact</a></li>
            <li>
              <a href="https://github.com/5pirit5eal/swim-rag" target="_blank" rel="noopener"
                >GitHub</a
              >
            </li>
          </ul>
        </div>
      </div>

      <div class="footer-bottom">
        <p>&copy; 2025 Swim RAG. Built with Vue.js and TypeScript.</p>
      </div>
    </div>
  </footer>
  <BaseModal :show="showPrivacyPolicyModal" @close="showPrivacyPolicyModal = false">
    <template #header><h1>Privacy Policy</h1></template>
    <template #body>
      <PrivacyPolicy />
    </template>
  </BaseModal>
</template>

<style scoped>
.app-footer {
  background-color: var(--color-background-soft, #f8f9fa);
  border-top: 1px solid var(--color-border, #e2e8f0);
  margin-top: auto;
  padding: 2rem 0 1rem;
}

.footer-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

.footer-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.footer-section h3,
.footer-section h4 {
  margin: 0 0 1rem 0;
  color: var(--color-heading, #2c3e50);
}

.footer-section h3 {
  font-size: 1.25rem;
}

.footer-section h4 {
  font-size: 1rem;
}

.footer-section p {
  margin: 0;
  color: var(--color-text-light, #64748b);
  line-height: 1.6;
}

.footer-links {
  list-style: none;
  padding: 0;
  margin: 0;
}

.footer-links li {
  margin-bottom: 0.5rem;
}

.footer-links a {
  color: var(--color-text, #374151);
  text-decoration: none;
  transition: color 0.2s;
}

.footer-links a:hover {
  color: var(--color-text-light, #64748b);
  text-decoration: underline;
}

.footer-bottom {
  padding-top: 1rem;
  border-top: 1px solid var(--color-border, #e2e8f0);
  text-align: center;
}

.footer-bottom p {
  margin: 0;
  color: var(--color-text-light, #64748b);
  font-size: 0.875rem;
}
</style>
