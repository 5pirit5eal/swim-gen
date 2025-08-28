<template>
  <div class="privacy-policy">
    <section class="intro">
      <p class="last-updated">
        <strong>{{ $t('privacy_policy.last_updated') }}</strong>
      </p>
    </section>

    <section>
      <h3>{{ $t('privacy_policy.data_protection_at_a_glance_title') }}</h3>
      <h4>{{ $t('privacy_policy.general_notes_title') }}</h4>
      <p>{{ $t('privacy_policy.general_notes_text') }}</p>

      <h4>{{ $t('privacy_policy.data_collection_on_this_website_title') }}</h4>
      <h5>{{ $t('privacy_policy.who_is_responsible_title') }}</h5>
      <p>{{ $t('privacy_policy.who_is_responsible_text') }}</p>

      <h5>{{ $t('privacy_policy.how_we_collect_data_title') }}</h5>
      <p>{{ $t('privacy_policy.how_we_collect_data_text') }}</p>

      <h5>{{ $t('privacy_policy.what_we_use_data_for_title') }}</h5>
      <p>{{ $t('privacy_policy.what_we_use_data_for_text') }}</p>

      <h5>{{ $t('privacy_policy.your_rights_title') }}</h5>
      <p>{{ $t('privacy_policy.your_rights_text') }}</p>
    </section>

    <section>
      <h3>{{ $t('privacy_policy.general_and_mandatory_info_title') }}</h3>
      <h4>{{ $t('privacy_policy.data_protection_title') }}</h4>
      <p>{{ $t('privacy_policy.data_protection_text') }}</p>

      <h4>{{ $t('privacy_policy.responsible_party_title') }}</h4>
      <p>{{ $t('privacy_policy.responsible_party_text') }}</p>
      <p v-html="$t('privacy_policy.responsible_party_details').replace(/\n/g, '<br>')"></p>
      <p>{{ $t('privacy_policy.responsible_party_definition') }}</p>

      <h4>{{ $t('privacy_policy.storage_duration_title') }}</h4>
      <p>{{ $t('privacy_policy.storage_duration_text') }}</p>

      <h4>{{ $t('privacy_policy.legal_basis_title') }}</h4>
      <p>{{ $t('privacy_policy.legal_basis_text') }}</p>

      <h4>{{ $t('privacy_policy.revocation_title') }}</h4>
      <p>{{ $t('privacy_policy.revocation_text') }}</p>

      <h4>{{ $t('privacy_policy.objection_right_title') }}</h4>
      <p>{{ $t('privacy_policy.objection_right_text_1') }}</p>
      <p>{{ $t('privacy_policy.objection_right_text_2') }}</p>

      <h4>{{ $t('privacy_policy.complaint_right_title') }}</h4>
      <p>{{ $t('privacy_policy.complaint_right_text') }}</p>

      <h4>{{ $t('privacy_policy.data_portability_right_title') }}</h4>
      <p>{{ $t('privacy_policy.data_portability_right_text') }}</p>

      <h4>{{ $t('privacy_policy.information_correction_deletion_title') }}</h4>
      <p>{{ $t('privacy_policy.information_correction_deletion_text') }}</p>

      <h4>{{ $t('privacy_policy.processing_restriction_right_title') }}</h4>
      <p>{{ $t('privacy_policy.processing_restriction_right_text_1') }}</p>
      <ul>
        <li v-for="item in $t('privacy_policy.processing_restriction_right_list').split('\n')" :key="item">{{ item }}
        </li>
      </ul>
      <p>{{ $t('privacy_policy.processing_restriction_right_text_2') }}</p>

      <h4>{{ $t('privacy_policy.ssl_tls_encryption_title') }}</h4>
      <p>{{ $t('privacy_policy.ssl_tls_encryption_text') }}</p>
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
  margin-top: 1.5rem;
  margin-bottom: 1rem;
  font-size: 1.5rem;
  font-weight: 600;
}

.privacy-policy h4 {
  color: var(--color-heading, #2c3e50);
  margin-top: 1.25rem;
  margin-bottom: 0.75rem;
  font-size: 1.25rem;
  font-weight: 600;
}

.privacy-policy h5 {
  color: var(--color-heading, #2c3e50);
  margin-top: 1rem;
  margin-bottom: 0.5rem;
  font-size: 1.1rem;
  font-weight: 600;
}

.privacy-policy section {
  margin-bottom: 1.5rem;
}

.privacy-policy section+section {
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
    font-size: 1.2rem;
  }

  .privacy-policy h4 {
    font-size: 1.1rem;
  }

  .privacy-policy h5 {
    font-size: 1rem;
  }

  .privacy-policy .intro {
    padding: 1rem;
  }
}
</style>
