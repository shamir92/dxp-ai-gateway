import { PublicShell } from "@/components/public-shell";

function LegalSection({
  id,
  title,
  children,
}: {
  id: string;
  title: string;
  children: React.ReactNode;
}) {
  return (
    <section id={id} className="border-t border-line-warm py-10 first:border-t-0 first:pt-0">
      <h2 className="font-heading text-[22px] font-semibold text-ink">{title}</h2>
      <div className="mt-4 space-y-4 text-[13px] leading-relaxed text-muted-warm">
        {children}
      </div>
    </section>
  );
}

function Callout({ children }: { children: React.ReactNode }) {
  return (
    <p className="rounded-xl border border-line-warm bg-cream-card p-4 text-ink">
      {children}
    </p>
  );
}

const subProcessors = [
  {
    provider: "Model providers",
    purpose: "Running inference for the models you select",
    data: "Request payload (prompts/files) as needed to return a response — not retained by our gateway after the call",
  },
  {
    provider: "Hosting & CDN",
    purpose: "Serving the website, API, and console",
    data: "IP address, request and access logs",
  },
  {
    provider: "Payment processor",
    purpose: "Checkout, credit top-ups, and invoices",
    data: "Billing identity (name, email, company, tax ID) and payment details — we do not store full card numbers",
  },
  {
    provider: "Email / messaging",
    purpose: "Account notices, support replies, and WhatsApp contact",
    data: "Email address or phone number and message content you send",
  },
];

export default function PrivacyPage() {
  return (
    <PublicShell>
      <section className="bg-cream px-14 py-16">
        <div className="mx-auto max-w-[760px]">
          <p className="text-[10px] font-semibold uppercase tracking-wider text-flame">
            Legal
          </p>
          <h1 className="mt-3 font-heading text-[42px] font-semibold text-ink">
            Privacy Policy
          </h1>
          <p className="mt-3 text-[13px] text-muted-warm">
            DXP AI Gateway · Last updated: 9 October 2026 · Includes rights under
            Indonesian Law No. 27 of 2022 (UU PDP), GDPR, and CCPA/CPRA where
            applicable
          </p>
        </div>
      </section>

      <section className="bg-white px-14 py-16">
        <div className="mx-auto max-w-[760px]">
          <p className="text-[13px] leading-relaxed text-muted-warm">
            This Privacy Policy explains how PT DXP Global Inovasi (“DXP”, “we”,
            “us”, or “our”) collects, uses, shares, and protects your information
            when you use DXP AI Gateway — our website, platform, API, and related
            services (the “Service”). By using the Service, you agree to the
            practices described here.
          </p>

          <LegalSection id="who" title="1. Who we are">
            <p>
              DXP AI Gateway is operated by PT DXP Global Inovasi, incorporated
              in the Republic of Indonesia, registered at Jl. Tarumajaya Raya,
              Sagara Makmur, Kec. Tarumajaya, Kabupaten Bekasi, Jawa Barat
              17211, Indonesia.
            </p>
            <p>
              For the purposes of the UU PDP (and where applicable, the GDPR), we
              act as the controller of personal data described in this policy.
              Contact:{" "}
              <a href="mailto:hello@fromdxp.com" className="text-flame">
                hello@fromdxp.com
              </a>
              .
            </p>
          </LegalSection>

          <LegalSection id="no-store" title="2. We do not store your content">
            <Callout>
              <strong>Important:</strong> We do not store prompts, completions,
              messages, files, or other request/response content that passes
              through the gateway. Content is processed in transit to fulfill the
              API call and is not retained on our systems after the response is
              returned (except as strictly required for short-lived operational
              buffering or retries, or when you explicitly opt in to a feature
              that requires retention).
            </Callout>
            <p>
              We do not use your content to train shared, general-purpose models.
              We do not sell your content or personal information.
            </p>
          </LegalSection>

          <LegalSection id="usage-logs" title="3. Token usage logs for billing">
            <p>
              To bill accurately and show spend in your dashboard, we log{" "}
              <strong>usage metadata only</strong>, including:
            </p>
            <ul className="list-disc space-y-2 pl-5">
              <li>Timestamp of the request</li>
              <li>API key identifier and account / organization ID</li>
              <li>Model name and endpoint</li>
              <li>
                Token counts (input, output, and cache where applicable) and
                resulting cost
              </li>
              <li>Request status and latency metrics</li>
            </ul>
            <p>
              These records power billing, credit consumption, usage analytics,
              abuse prevention, and support. They are not used to reconstruct or
              store the content of your prompts or model outputs.
            </p>
          </LegalSection>

          <LegalSection id="collect" title="4. Information we collect">
            <p>
              We collect information needed to run the Service. Categories
              include:
            </p>
            <ul className="list-disc space-y-2 pl-5">
              <li>
                <strong>Account & identity.</strong> Email address and, if you
                sign in with a third-party identity provider, basic profile
                information (such as your name and email) used to create and
                authenticate your account.
              </li>
              <li>
                <strong>Billing & profile.</strong> Account type (individual or
                business), full name, company and contact name, phone number,
                postal address, and tax identification number (optional).
              </li>
              <li>
                <strong>Payments.</strong> Payments are processed by our payment
                provider. We do not collect or store full payment card numbers.
                We receive limited transaction details such as order ID, product,
                amount, and status.
              </li>
              <li>
                <strong>Credit & usage records.</strong> Prepaid credit balance
                and a history of top-ups and metered usage events (see Section 3).
              </li>
              <li>
                <strong>Technical & log data.</strong> When you access the
                Service, standard technical information such as IP address,
                browser and device type, and timestamps, used to serve and secure
                the Service.
              </li>
              <li>
                <strong>Communications.</strong> If you email us at
                hello@fromdxp.com or message us on WhatsApp, we receive the
                contact details and content you choose to send.
              </li>
            </ul>
            <p>
              We do not intentionally collect specific (sensitive) personal data
              under the UU PDP (such as health, biometric, genetic, or children’s
              data). Please do not send such data through the API or by email.
            </p>
          </LegalSection>

          <LegalSection id="cookies" title="5. Cookies and local storage">
            <p>
              We use your browser’s local storage to keep you signed in (your
              authentication session). This is strictly necessary for the Service
              to function.
            </p>
            <p>
              We avoid analytics and advertising cookies where possible. If we
              introduce non-essential cookies or similar technologies, we will
              update this policy first and ask for your consent before they are
              placed on your device, where required. Rejecting optional cookies
              will not prevent you from using the site.
            </p>
            <p>
              Some third-party resources we load (such as fonts or icon assets)
              may be served by a CDN that receives your IP address as part of
              delivering those files.
            </p>
          </LegalSection>

          <LegalSection id="why" title="6. How we use your information">
            <p>
              Where the GDPR applies, the corresponding legal basis is shown in
              parentheses.
            </p>
            <ul className="list-disc space-y-2 pl-5">
              <li>
                <strong>Provide, operate, and maintain the Service</strong> —
                authenticate API keys, route requests, meter usage, and manage
                credits (performance of contract).
              </li>
              <li>
                <strong>Secure your account</strong> — prevent fraud and abuse,
                debug issues, keep the gateway reliable (legitimate interests).
              </li>
              <li>
                <strong>Process payments, issue invoices, and meet tax
                obligations</strong> (performance of contract; legal obligation).
              </li>
              <li>
                <strong>Support</strong> — respond to questions via email or
                WhatsApp (legitimate interests / steps prior to contract).
              </li>
              <li>
                <strong>Comply with applicable laws</strong> and respond to lawful
                requests (legal obligation).
              </li>
            </ul>
            <p>
              We do not use personal data for automated decision-making that
              produces legal or similarly significant effects for you, and we do
              not sell personal information or share it for cross-context
              behavioral advertising as those terms are defined under California
              law.
            </p>
          </LegalSection>

          <LegalSection id="share" title="7. How we share your information">
            <p>
              We share personal information only with service providers
              (sub-processors) that help us run the Service, and only as needed
              for the purposes below. We do not sell your personal information.
            </p>
            <div className="overflow-x-auto rounded-xl border border-line-warm">
              <table className="w-full min-w-[520px] text-left text-[12px]">
                <thead className="bg-cream-card text-ink">
                  <tr>
                    <th className="px-4 py-3 font-semibold">Provider</th>
                    <th className="px-4 py-3 font-semibold">Purpose</th>
                    <th className="px-4 py-3 font-semibold">Data shared</th>
                  </tr>
                </thead>
                <tbody>
                  {subProcessors.map((row) => (
                    <tr key={row.provider} className="border-t border-line-warm">
                      <td className="px-4 py-3 text-ink">{row.provider}</td>
                      <td className="px-4 py-3">{row.purpose}</td>
                      <td className="px-4 py-3">{row.data}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            <p>
              We may also disclose information if required by law, to enforce our
              terms, or to protect the rights, property, or safety of our users or
              others. If we are involved in a merger, acquisition, or asset sale,
              personal information may be transferred, and we will notify you of
              any change in control or use of your information as required by law.
            </p>
          </LegalSection>

          <LegalSection id="transfers" title="8. International data transfers">
            <p>
              Because we and our providers may operate globally, your information
              may be processed in countries other than the one in which you live,
              including countries that may not provide the same level of data
              protection. Where required, we rely on appropriate safeguards (such
              as standard contractual clauses or other lawful transfer mechanisms
              under the UU PDP and GDPR).
            </p>
          </LegalSection>

          <LegalSection id="retention" title="9. Data retention">
            <p>
              <strong>Request/response content:</strong> not stored after the
              request completes (see Section 2).
            </p>
            <p>
              <strong>Usage and billing logs:</strong> kept for as long as needed
              for billing, accounting, dispute resolution, and legal compliance
              (typically the current billing period plus a reasonable statutory
              window).
            </p>
            <p>
              <strong>Account and profile data:</strong> kept for as long as your
              account is active. Transaction, invoice, and tax records are
              retained for the period required by applicable accounting and tax
              laws.
            </p>
            <p>
              When information is no longer needed, we delete or irreversibly
              anonymize it. You may request deletion of your account at any time
              (see Section 10).
            </p>
          </LegalSection>

          <LegalSection id="security" title="10. Security">
            <p>
              We protect your data with measures including encryption in transit
              (HTTPS), access controls limiting who can read service data, and
              tightly scoped access to payment and service secrets. No method of
              transmission or storage is completely secure, but we work to
              protect your information and to address vulnerabilities promptly.
            </p>
            <p>
              If a personal data breach occurs that affects you, we will notify
              you and the relevant authority as required by the UU PDP (no later
              than 3 × 24 hours where applicable) and by other laws where they
              apply.
            </p>
          </LegalSection>

          <LegalSection id="rights" title="11. Your rights">
            <p>
              Depending on where you live, you have rights over your personal
              information. Subject to applicable law, these generally include the
              right to access, correct, delete, or export your data; to object to
              or restrict certain processing; and to withdraw consent.
            </p>
            <ul className="list-disc space-y-2 pl-5">
              <li>
                <strong>Indonesia (UU PDP, Law No. 27 of 2022).</strong> You have
                the right to be informed, access, correct, erase, and obtain your
                personal data; to withdraw consent; to object to certain
                processing; and to seek redress for misuse of your data.
              </li>
              <li>
                <strong>EEA & UK (GDPR).</strong> In addition to the rights
                above, you may lodge a complaint with your local data protection
                supervisory authority.
              </li>
              <li>
                <strong>California (CCPA/CPRA).</strong> You have the right to
                know what personal information we collect, to delete it, to
                correct it, and to be free from discrimination for exercising
                your rights. We do not sell or share your personal information as
                those terms are defined under California law.
              </li>
            </ul>
            <p>
              To exercise any of these rights, email{" "}
              <a href="mailto:hello@fromdxp.com" className="text-flame">
                hello@fromdxp.com
              </a>{" "}
              (subject: “Data subject request”). We will respond within the
              timeframe required by applicable law and may need to verify your
              identity first. Many profile fields can be updated in account
              settings after you sign in.
            </p>
          </LegalSection>

          <LegalSection id="children" title="12. Children’s privacy">
            <p>
              The Service is not directed to children, and we do not knowingly
              collect personal information from anyone under the age of 16 (or
              the minimum age required in your jurisdiction). If you believe a
              child has provided us with personal information, contact{" "}
              <a href="mailto:hello@fromdxp.com" className="text-flame">
                hello@fromdxp.com
              </a>{" "}
              and we will delete it.
            </p>
          </LegalSection>

          <LegalSection id="third-party" title="13. Third-party links and services">
            <p>
              The Service relies on third-party providers (including model
              providers, payment, hosting, and messaging tools) and may link to
              third-party sites. Their handling of your information is governed
              by their own privacy policies, and we encourage you to review them.
            </p>
          </LegalSection>

          <LegalSection id="changes" title="14. Changes to this policy">
            <p>
              We may update this Privacy Policy from time to time. When we do, we
              will revise the “Last updated” date above, and for material changes
              we will provide additional notice where appropriate. Your continued
              use of the Service after an update means you accept the revised
              policy.
            </p>
          </LegalSection>

          <LegalSection id="contact" title="15. Contact us">
            <p>
              If you have questions about this Privacy Policy or how we handle
              your information, contact us at{" "}
              <a href="mailto:hello@fromdxp.com" className="text-flame">
                hello@fromdxp.com
              </a>
              .
            </p>
            <p>
              PT DXP Global Inovasi
              <br />
              Jl. Tarumajaya Raya, Sagara Makmur, Kec. Tarumajaya, Kabupaten
              Bekasi, Jawa Barat 17211, Indonesia
            </p>
          </LegalSection>
        </div>
      </section>
    </PublicShell>
  );
}
