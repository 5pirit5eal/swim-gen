<template>
  <div class="impressum-container">
    <section>
      <h3>{{ $t('impressum.title') }}</h3>
      <p>
        {{ $t('impressum.name') }}<br />
        {{ $t('impressum.address') }}<br />
        {{ $t('impressum.city') }}
      </p>
    </section>
    <section>
      <h3>{{ $t('impressum.contact_title') }}</h3>
      <p>
        {{ $t('impressum.phone') }}<br />
        {{ $t('impressum.email') }}
      </p>
    </section>
    <section>
      <p>
        <small>{{ $t('impressum.source') }}</small>
      </p>
    </section>
  </div>
</template>

<style scoped>
.impressum-container {
  padding: 1rem;
  line-height: 1.6;
  color: var(--color-text);
}

h3 {
  color: var(--color-heading);
  margin-top: 1.5rem;
  margin-bottom: 1rem;
  font-size: 1.25rem;
  font-weight: 600;
}

p {
  margin-bottom: 1rem;
}

small {
  font-size: 0.8rem;
  color: var(--color-text-light);
}
</style>
