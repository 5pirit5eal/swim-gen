model            = "gemini-2.5-flash"
small_model      = "gemini-2.5-flash-lite"
embedding_name   = "embeddings"
embedding_model  = "gemini-embedding-001"
embedding_size   = 768
log_level        = "INFO"
domain_url       = "swim-gen.com"
outputs_location = "../0-config"
resend_dns_records = [
  {
    name  = "resend._domainkey"
    value = "p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDaFJpmOEXCoT7P4TUaSzIObPkac3MfjxTOHzLct7wTupWdAodZ+btJu6K9av427giwx/OqsZhHHToM0dgEI/StC/D/eiN4QS1Ldc8fjPOe7LKmnbCHx208Dm/CsO5Wetcne3VW1E7SGlx1ApNsUpRB6JoMzZQ1uRKUSt3Yl4mSOQIDAQAB"
    type  = "TXT"
    ttl   = 300
  },
  {
    name  = "send"
    value = "10 feedback-smtp.eu-west-1.amazonses.com"
    type  = "MX"
    ttl   = 300
  },
  {
    name  = "send"
    value = "\"v=spf1 include:amazonses.com ~all\""
    type  = "TXT"
    ttl   = 300
  },
  {
    name  = "_dmarc"
    value = "\"v=DMARC1; p=none;\""
    type  = "TXT"
    ttl   = 300
  }
]
