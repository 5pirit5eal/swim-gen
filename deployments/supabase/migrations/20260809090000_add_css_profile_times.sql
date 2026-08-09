alter table profiles
  add column css_200m_seconds integer,
  add column css_400m_seconds integer;

alter table profiles
  add constraint profiles_css_200m_positive
    check (css_200m_seconds is null or css_200m_seconds > 0),
  add constraint profiles_css_400m_positive
    check (css_400m_seconds is null or css_400m_seconds > 0),
  add constraint profiles_css_times_ordered
    check (
      css_200m_seconds is null
      or css_400m_seconds is null
      or css_400m_seconds > css_200m_seconds
    );
