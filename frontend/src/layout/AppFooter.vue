<script setup lang="ts">
import IconHeart from '@/components/icons/IconHeart.vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
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
              <router-link to="/privacy">{{ t('footer.privacy_policy') }}</router-link>
            </li>
            <li>
              <router-link to="/terms">{{ t('footer.terms_of_service') }}</router-link>
            </li>
            <li>
              <router-link to="/impressum">{{ t('footer.impressum') }}</router-link>
            </li>
          </ul>
        </div>

        <div class="footer-section">
          <h4>{{ t('footer.about') }}</h4>
          <ul class="footer-links">
            <li>
              <a href="https://github.com/5pirit5eal/swim-gen" target="_blank" rel="noopener">{{
                t('footer.github')
              }}</a>
            </li>
            <li>
              <a href="https://github.com/5pirit5eal/swim-gen" target="_blank" rel="noopener">{{
                t('footer.roadmap')
              }}</a>
            </li>
            <li>
              <a href="https://buymeacoffee.com/spiritseal" target="_blank" rel="noopener">
                <IconHeart class="icon-heart" />
                {{ t('footer.sponsor') }}
              </a>
            </li>
          </ul>
        </div>
      </div>

      <div class="footer-bottom">
        <p>{{ t('footer.copyright') }}</p>
      </div>
    </div>
  </footer>
</template>

<style scoped>
.app-footer {
  background-color: var(--color-background-soft);
  border-top: 1px solid var(--color-border);
  margin-top: auto;
  padding: 2rem 0 1rem 0;
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

@media (max-width: 740px) {
  .footer-content {
    justify-items: left;
  }

  .footer-section {
    margin-bottom: 1.5rem;
  }
}

.footer-section h3,
.footer-section h4 {
  margin: 0 0 1rem 0;
  color: var(--color-heading);
}

.footer-section h3 {
  font-size: 1.25rem;
}

.footer-section h4 {
  font-size: 1rem;
}

.footer-section p {
  margin: 0;
  color: var(--color-text-light);
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
  color: var(--color-text);
  text-decoration: none;
  transition: color 0.2s;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.footer-links a:hover {
  color: var(--color-text-light);
  text-decoration: underline;
}

.icon-heart {
  width: 1rem;
  height: 1rem;
  color: #732441;
  transition: color 0.2s;
}

.footer-links a:hover .icon-heart {
  color: var(--color-text-light);
}

.footer-bottom {
  padding-top: 1rem;
  border-top: 1px solid var(--color-border);
  text-align: center;
}

.footer-bottom p {
  margin: 0;
  color: var(--color-text-light);
  font-size: 0.875rem;
}
</style>
