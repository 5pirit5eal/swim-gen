<script setup lang="ts">
const VITE_IMPRESSUM_EMAIL = import.meta.env.VITE_IMPRESSUM_EMAIL
</script>

<template>
  <div class="legal-page">
    <div class="container">
      <div class="terms-of-service">
        <h1>{{ $t('terms_of_service.title') }}</h1>
        <section class="intro">
          <p class="last-updated">
            <strong>{{ $t('terms_of_service.last_updated') }}</strong>
          </p>
          <p>{{ $t('terms_of_service.read_carefully') }}</p>
          <p>{{ $t('terms_of_service.provided_by') }}</p>
          <p>{{ $t('terms_of_service.owner') }}<br /></p>
          <p>{{ $t('terms_of_service.owner_contact', { mail: VITE_IMPRESSUM_EMAIL }) }}</p>
        </section>

        <section>
          <h3>{{ $t('terms_of_service.at_a_glance_title') }}</h3>
          <p>{{ $t('terms_of_service.at_a_glance') }}</p>
        </section>

        <section>
          <h3>{{ $t('terms_of_service.terms_of_use_title') }}</h3>
          <p>{{ $t('terms_of_service.terms_of_use') }}</p>
          <p>{{ $t('terms_of_service.additional_conditions') }}</p>
          <p>{{ $t('terms_of_service.requirements') }}</p>
          <ul>
            <li v-for="i in 1" :key="i">{{ $t(`terms_of_service.requirements_list.${i - 1}`) }}</li>
          </ul>
        </section>

        <section>
          <h3>{{ $t('terms_of_service.content_title') }}</h3>
          <p>{{ $t('terms_of_service.content') }}</p>
          <p>{{ $t('terms_of_service.content_effort') }}</p>
        </section>

        <section>
          <h3>{{ $t('terms_of_service.rights_title') }}</h3>
          <p>{{ $t('terms_of_service.rights') }}</p>
          <p>{{ $t('terms_of_service.rights_usage') }}</p>
          <p>{{ $t('terms_of_service.rights_prohibited') }}</p>
          <p>{{ $t('terms_of_service.rights_personal_use') }}</p>
          <p>{{ $t('terms_of_service.rights_limitation') }}</p>
        </section>

        <section>
          <h3>{{ $t('terms_of_service.user_content_title') }}</h3>
          <p>{{ $t('terms_of_service.user_content') }}</p>
          <p>{{ $t('terms_of_service.user_content_confirm') }}</p>
        </section>

        <section>
          <h3>{{ $t('terms_of_service.external_resources_title') }}</h3>
          <p>{{ $t('terms_of_service.external_resources') }}</p>
          <p>{{ $t('terms_of_service.external_resources_conditions') }}</p>
        </section>

        <section>
          <h3>{{ $t('terms_of_service.acceptable_use_title') }}</h3>
          <p>{{ $t('terms_of_service.acceptable_use') }}</p>
          <p>{{ $t('terms_of_service.acceptable_use_responsibility') }}</p>
          <p>{{ $t('terms_of_service.acceptable_use_rights') }}</p>
          <ul>
            <li v-for="i in 4" :key="i">
              {{ $t(`terms_of_service.acceptable_use_list.${i - 1}`) }}
            </li>
          </ul>
        </section>

        <section>
          <h3>{{ $t('terms_of_service.privacy_policy_title') }}</h3>
          <p>{{ $t('terms_of_service.privacy_policy') }}</p>
        </section>

        <section>
          <h3>{{ $t('terms_of_service.ip_rights_title') }}</h3>
          <p>{{ $t('terms_of_service.ip_rights') }}</p>
          <p>{{ $t('terms_of_service.ip_rights_trademarks') }}</p>
        </section>

        <section>
          <h3>{{ $t('terms_of_service.changes_title') }}</h3>
          <p>{{ $t('terms_of_service.changes') }}</p>
          <p>{{ $t('terms_of_service.changes_effect') }}</p>
          <p>{{ $t('terms_of_service.changes_acceptance') }}</p>
          <p>{{ $t('terms_of_service.changes_previous') }}</p>
          <p>{{ $t('terms_of_service.changes_law') }}</p>
        </section>

        <section>
          <h3>{{ $t('terms_of_service.assignment_title') }}</h3>
          <p>{{ $t('terms_of_service.assignment') }}</p>
          <p>{{ $t('terms_of_service.assignment_user') }}</p>
        </section>

        <section>
          <h3>{{ $t('terms_of_service.contacts_title') }}</h3>
          <p>{{ $t('terms_of_service.contacts') }}</p>
        </section>

        <section>
          <h3>{{ $t('terms_of_service.severability_title') }}</h3>
          <p>{{ $t('terms_of_service.severability') }}</p>
        </section>

        <section>
          <h3>{{ $t('terms_of_service.governing_law_title') }}</h3>
          <p>{{ $t('terms_of_service.governing_law') }}</p>
        </section>

        <section>
          <h3>{{ $t('terms_of_service.venue_title') }}</h3>
          <p>{{ $t('terms_of_service.venue') }}</p>
        </section>

        <section>
          <h3>{{ $t('terms_of_service.dispute_title') }}</h3>
          <h3>{{ $t('terms_of_service.amicable_title') }}</h3>
          <p>{{ $t('terms_of_service.amicable') }}</p>
          <p>{{ $t('terms_of_service.amicable_right') }}</p>
          <p>{{ $t('terms_of_service.amicable_complaint') }}</p>
          <p>{{ $t('terms_of_service.amicable_process') }}</p>
        </section>
      </div>
    </div>
  </div>
</template>

<style scoped>
.legal-page {
  padding: 0;
}

.container {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 1rem;
}

.terms-of-service {
  max-width: 800px;
  margin: 2rem auto;
  padding: 2rem;
  line-height: 1.6;
  color: var(--color-text);
  background-color: var(--color-background-soft);
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

.terms-of-service h1 {
  text-align: center;
  margin-bottom: 2rem;
  color: var(--color-heading);
}

.terms-of-service h3 {
  color: var(--color-heading, #2c3e50);
  margin-top: 1rem;
  margin-bottom: 1rem;
  font-size: 1.25rem;
  font-weight: 600;
}

.terms-of-service section {
  margin-bottom: 1.5rem;
}

.terms-of-service section+section {
  margin-top: 1.5rem;
}

.terms-of-service p {
  margin-bottom: 1rem;
}

.terms-of-service ul {
  margin: 1rem 0;
  padding-left: 1.5rem;
}

.terms-of-service li {
  margin-bottom: 0.5rem;
}

.terms-of-service .last-updated {
  font-size: 0.9rem;
  color: var(--color-text-light, #64748b);
  background: var(--color-background-soft, #f8f9fa);
  padding: 0.75rem;
  border-radius: 0.375rem;
  border-left: 3px solid var(--color-border, #e2e8f0);
}

.terms-of-service .intro {
  display: block;
  background: var(--color-background);
  padding: 1.5rem;
  border-radius: 0.5rem;
  margin-bottom: 2rem;
}

.terms-of-service a {
  color: var(--color-text, #374151);
  text-decoration: underline;
}

.terms-of-service a:hover {
  color: var(--color-text-light, #64748b);
}

.terms-of-service strong {
  font-weight: 600;
  color: var(--color-heading, #2c3e50);
}

/* Mobile responsiveness */
@media (max-width: 768px) {
  .terms-of-service {
    max-width: 100%;
  }

  .terms-of-service h3 {
    font-size: 1.1rem;
  }

  .terms-of-service .intro {
    padding: 1rem;
  }
}
</style>
