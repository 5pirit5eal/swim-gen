<script setup lang="ts">
import BaseModal from '@/components/ui/BaseModal.vue'
import PrivacyPolicy from '../legal/PrivacyPolicy.vue'
import TermsOfService from '../legal/TermsOfService.vue'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

// Footer component with privacy links and legal information
// Handlers for modal dialogs will be implemented when needed

const showPrivacyPolicyModal = ref(false)
const showTermsOfServiceModal = ref(false)

const { t } = useI18n()

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
          <h3>{{ t('app.name') }}</h3>
          <p>{{ t('app.slogan') }}</p>
        </div>

        <div class="footer-section">
          <h4>{{ t('footer.privacy_legal') }}</h4>
          <ul class="footer-links">
            <li>
              <a href="#privacy" id="privacy-policy" @click="showPrivacyPolicyModal = true">{{
                t('footer.privacy_policy')
              }}</a>
            </li>
            <li>
              <a href="#terms" @click="showTermsOfServiceModal = true">{{
                t('footer.terms_of_service')
              }}</a>
            </li>
            <li>
              <a href="#impressum" @click.prevent="showImpressum">{{ t('footer.impressum') }}</a>
            </li>
          </ul>
        </div>

        <div class="footer-section">
          <h4>{{ t('footer.about') }}</h4>
          <ul class="footer-links">
            <li>
              <a href="#contact" @click.prevent="showContact">{{ t('footer.contact') }}</a>
            </li>
            <li>
              <a href="https://github.com/5pirit5eal/swim-rag" target="_blank" rel="noopener">{{
                t('footer.github')
              }}</a>
            </li>
          </ul>
        </div>
      </div>

      <div class="footer-bottom">
        <p>{{ t('footer.copyright') }}</p>
      </div>
    </div>
  </footer>
  <BaseModal :show="showPrivacyPolicyModal" @close="showPrivacyPolicyModal = false">
    <template #header>
      <h1>{{ t('footer.privacy_policy') }}</h1>
    </template>
    <template #body>
      <PrivacyPolicy />
    </template>
  </BaseModal>

  <BaseModal :show="showTermsOfServiceModal" @close="showTermsOfServiceModal = false">
    <template #header>
      <h1>{{ t('footer.terms_of_service') }}</h1>
    </template>
    <template #body>
      <TermsOfService />
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
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 1rem;
}

.footer-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  justify-items: center;
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
