<template>
  <div class="privacy-policy">
    <section class="intro">
      <p class="last-updated">
        <strong>{{ $t('privacy_policy.last_updated') }}</strong>
      </p>
      <p>{{ $t('privacy_policy.intro') }}</p>
    </section>

    <section>
      <h3>{{ $t('privacy_policy.privacy_by_design_title') }}</h3>
      <p>{{ $t('privacy_policy.privacy_by_design') }}</p>
    </section>

    <section>
      <h3>{{ $t('privacy_policy.info_we_collect_title') }}</h3>
      <p>{{ $t('privacy_policy.info_we_collect') }}</p>
    </section>

    <section>
      <h3>{{ $t('privacy_policy.technical_info_title') }}</h3>
      <p>{{ $t('privacy_policy.technical_info') }}</p>
    </section>

    <section>
      <h3>{{ $t('privacy_policy.cookies_title') }}</h3>
      <p>{{ $t('privacy_policy.cookies') }}</p>
    </section>

    <section>
      <h3>{{ $t('privacy_policy.future_data_donation_title') }}</h3>
      <p>{{ $t('privacy_policy.future_data_donation') }}</p>
    </section>

    <section>
      <h3>{{ $t('privacy_policy.third_party_title') }}</h3>
      <p>{{ $t('privacy_policy.third_party') }}</p>
    </section>

    <section>
      <h3>{{ $t('privacy_policy.gdpr_rights_title') }}</h3>
      <p>{{ $t('privacy_policy.gdpr_rights') }}</p>
    </section>

    <section>
      <h3>{{ $t('privacy_policy.data_security_title') }}</h3>
      <p>{{ $t('privacy_policy.data_security') }}</p>
    </section>

    <section>
      <h3>{{ $t('privacy_policy.intl_data_transfers_title') }}</h3>
      <p>{{ $t('privacy_policy.intl_data_transfers') }}</p>
    </section>

    <section>
      <h3>{{ $t('privacy_policy.children_title') }}</h3>
      <p>{{ $t('privacy_policy.children') }}</p>
    </section>

    <section>
      <h3>{{ $t('privacy_policy.changes_title') }}</h3>
      <p>{{ $t('privacy_policy.changes') }}</p>
    </section>

    <section class="contact">
      <h3>{{ $t('privacy_policy.contact_title') }}</h3>
      <p>{{ $t('privacy_policy.contact') }}</p>
    </section>

    <section class="commitment">
      <h3>{{ $t('privacy_policy.commitment_title') }}</h3>
      <p>{{ $t('privacy_policy.commitment') }}</p>
    </section>
  </div>
</template>

<style scoped>
.privacy-policy {
  margin-right: 0.5rem;
  line-height: 1.6;
  color: var(--color-text);
}

.privacy-policy h3 {
  color: var(--color-heading, #2c3e50);
  margin-top: 1rem;
  margin-bottom: 1rem;
  font-size: 1.25rem;
  font-weight: 600;
}

.privacy-policy section {
  margin-bottom: 1.5rem;
}

.privacy-policy section + section {
  margin-top: 1.5rem;
}

.privacy-policy p {
  margin-bottom: 1rem;
}

.privacy-policy ul {
  margin: 1rem 0;
  padding-left: 1.5rem;
}

.privacy-policy li {
  margin-bottom: 0.5rem;
}

.privacy-policy .last-updated {
  font-size: 0.9rem;
  color: var(--color-text-light, #64748b);
  background: var(--color-background-soft, #f8f9fa);
  padding: 0.75rem;
  border-radius: 0.375rem;
  border-left: 3px solid var(--color-border, #e2e8f0);
}

.privacy-policy .intro {
  background: var(--color-background-soft, #f8f9fa);
  padding: 1.5rem;
  border-radius: 0.5rem;
  margin-bottom: 2rem;
}

.privacy-policy .contact {
  background: var(--color-background-soft, #f8f9fa);
  padding: 1.5rem;
  border-radius: 0.5rem;
  margin-top: 2rem;
}

.privacy-policy .commitment {
  background: linear-gradient(
    135deg,
    var(--color-background-soft, #f8f9fa) 0%,
    var(--color-background, #ffffff) 100%
  );
  padding: 1.5rem;
  border-radius: 0.5rem;
  margin-top: 2rem;
  border: 1px solid var(--color-border, #e2e8f0);
}

.privacy-policy a {
  color: var(--color-text, #374151);
  text-decoration: underline;
}

.privacy-policy a:hover {
  color: var(--color-text-light, #64748b);
}

.privacy-policy strong {
  font-weight: 600;
  color: var(--color-heading, #2c3e50);
}

/* Mobile responsiveness */
@media (max-width: 768px) {
  .privacy-policy {
    max-width: 100%;
  }

  .privacy-policy h3 {
    font-size: 1.1rem;
  }

  .privacy-policy .intro,
  .privacy-policy .contact,
  .privacy-policy .commitment {
    padding: 1rem;
  }
}
</style>
