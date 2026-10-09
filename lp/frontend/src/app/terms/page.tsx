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

export default function TermsPage() {
  return (
    <PublicShell>
      <section className="bg-cream px-14 py-16">
        <div className="mx-auto max-w-[760px]">
          <p className="text-[10px] font-semibold uppercase tracking-wider text-flame">
            Legal
          </p>
          <h1 className="mt-3 font-heading text-[42px] font-semibold text-ink">
            Terms of Use
          </h1>
          <p className="mt-3 text-[13px] text-muted-warm">
            DXP AI Gateway · Last updated: 9 October 2026
          </p>
        </div>
      </section>

      <section className="bg-white px-14 py-16">
        <div className="mx-auto max-w-[760px]">
          <p className="text-[13px] leading-relaxed text-muted-warm">
            These Terms of Use (“Terms”) govern your access to and use of DXP AI
            Gateway’s website, platform, API, and services (the “Service”).
            “DXP”, “we”, “us”, and “our” refer to PT DXP Global Inovasi, the
            operator of the Service; “you” refers to the person or entity using
            it. By creating an account or using the Service, you agree to these
            Terms and our Privacy Policy.
          </p>
          <p className="text-[13px] leading-relaxed text-muted-warm">
            If you are using the Service on behalf of an organization, you
            represent that you have authority to bind that organization to these
            Terms, and “you” includes that organization.
          </p>

          <LegalSection id="service" title="1. The Service">
            <p>
              DXP AI Gateway provides access to hosted AI models through
              OpenAI-compatible and Anthropic-compatible APIs, including usage
              analytics, key management, and prepaid credit-based billing. The
              Service is offered “as available” and may evolve as we add, modify,
              or remove features.
            </p>
          </LegalSection>

          <LegalSection id="eligibility" title="2. Eligibility and your account">
            <p>
              To use the Service you must be at least 16 years old (or the
              minimum age required in your jurisdiction) and able to enter into a
              binding contract. The Service is not directed to children below
              that age.
            </p>
            <p>
              You are responsible for keeping your account credentials and API
              keys secure, for all activity that occurs under your account or
              keys, and for the accuracy of the information you provide
              (including billing details). Notify us immediately at{" "}
              <a href="mailto:hello@fromdxp.com" className="text-flame">
                hello@fromdxp.com
              </a>{" "}
              if you suspect unauthorized use of an API key.
            </p>
            <p>
              You may close your account at any time. We may suspend or terminate
              accounts as described in Section 9.
            </p>
          </LegalSection>

          <LegalSection id="billing" title="3. Credits, billing, and taxes">
            <p>The Service is billed on a prepaid credit and usage basis.</p>
            <ul className="list-disc space-y-2 pl-5">
              <li>
                <strong>Purchasing credits.</strong> You may purchase credit
                packs or subscribe to a plan at the prices shown in the Service.
              </li>
              <li>
                <strong>Payment.</strong> Payments are processed by our payment
                providers. By purchasing, you agree to their applicable terms and
                authorize them to charge your selected payment method. We do not
                store your full payment card details.
              </li>
              <li>
                <strong>Taxes.</strong> Prices may exclude applicable taxes (such
                as VAT or GST). Final totals are calculated at checkout. You are
                responsible for taxes that are your obligation under applicable
                law.
              </li>
              <li>
                <strong>Consumption.</strong> Credits are consumed when you use
                metered features (including model token usage). Once consumed,
                credits cannot be restored.
              </li>
              <li>
                <strong>Refunds.</strong> Except where required by law, purchased
                credits are non-refundable and have no cash value. Unused credits
                may be forfeited if your account is terminated for a breach of
                these Terms.
              </li>
              <li>
                <strong>Pricing changes.</strong> We may change prices, pack
                sizes, or consumption rates. Changes apply to future purchases
                and consumption; credits already purchased remain valid at the
                rate in effect when consumed, subject to any expiration we
                communicate in advance.
              </li>
            </ul>
          </LegalSection>

          <LegalSection id="acceptable-use" title="4. Acceptable use">
            <p>You agree not to use the Service to:</p>
            <ul className="list-disc space-y-2 pl-5">
              <li>
                violate any law, regulation, or third party’s rights, including
                intellectual property, privacy, and publicity rights;
              </li>
              <li>
                generate, store, or distribute content that is unlawful, harmful,
                harassing, defamatory, hateful, sexually exploitative of minors,
                or that incites violence;
              </li>
              <li>
                create or distribute malware, attempt unauthorized access, probe
                or scan for vulnerabilities without permission, or interfere with
                the integrity or performance of the Service;
              </li>
              <li>
                reverse engineer, decompile, or otherwise attempt to extract the
                source code or underlying models of the Service, except to the
                extent that applicable law expressly permits;
              </li>
              <li>
                use automated means to scrape or extract data beyond documented
                interfaces, or to evade rate limits or access controls;
              </li>
              <li>
                resell, sublicense, or otherwise commercially exploit the Service
                in a way not permitted by these Terms;
              </li>
              <li>
                use the Service to make decisions that produce legal or similarly
                significant effects about a person without appropriate human
                review.
              </li>
            </ul>
            <p>
              We may investigate suspected violations and take any action we
              believe appropriate, including warning, throttling, suspending, or
              terminating your access.
            </p>
          </LegalSection>

          <LegalSection id="content" title="5. Your content; AI inputs and outputs">
            <p>
              “Inputs” are the data, prompts, files, and other content you submit
              to the Service. “Outputs” are the results the Service returns to
              you based on your Inputs.
            </p>
            <Callout>
              <strong>No content storage.</strong> We do not store your Inputs or
              Outputs (prompts, completions, messages, or other request/response
              content) in the gateway after a request is fulfilled. Content is
              processed in transit to deliver the model response.
            </Callout>
            <Callout>
              <strong>Token usage logs for billing.</strong> We record usage
              metadata only — including API key identifier, model, input/output
              token counts, cost, request status, and timing — for billing, usage
              dashboards, abuse prevention, and support. We do not use these logs
              to store or reconstruct your content.
            </Callout>
            <ul className="list-disc space-y-2 pl-5">
              <li>
                <strong>Ownership.</strong> As between you and us, you retain
                ownership of your Inputs and, subject to applicable law and these
                Terms, the Outputs generated for you. You are responsible for
                having the rights and permissions necessary to submit Inputs and
                use Outputs.
              </li>
              <li>
                <strong>License to us.</strong> You grant us a limited,
                non-exclusive license to process your Inputs and Outputs solely
                to operate, secure, and provide the Service to you. We will not
                use your Inputs to train shared, general-purpose models without
                your separate consent.
              </li>
              <li>
                <strong>Output disclaimer.</strong> AI Outputs may be inaccurate,
                incomplete, biased, or otherwise unsuitable for your purpose.
                Outputs do not constitute professional advice. You are
                responsible for reviewing Outputs before relying on them.
              </li>
              <li>
                <strong>No regulated decisions.</strong> Do not use the Service
                as the sole basis for decisions that have legal or similarly
                significant effects on individuals without appropriate human
                oversight and compliance with applicable law.
              </li>
            </ul>
            <p>
              See our Privacy Policy for how we handle account personal data.
            </p>
          </LegalSection>

          <LegalSection id="ip" title="6. Intellectual property">
            <p>
              The Service, including its software, designs, text, graphics, and
              other materials we provide, is owned by us or our licensors and is
              protected by intellectual property laws. We grant you a limited,
              non-exclusive, non-transferable, revocable license to access and
              use the Service in accordance with these Terms. No other rights are
              granted by implication or otherwise.
            </p>
            <p>
              You may submit feedback, suggestions, or ideas about the Service.
              If you do, you grant us a perpetual, irrevocable, royalty-free
              license to use them for any purpose, without obligation to you.
            </p>
          </LegalSection>

          <LegalSection id="third-party" title="7. Third-party services and privacy">
            <p>
              The Service relies on third-party services (including model
              providers, authentication, payment processing, and hosting). Your
              use of those services is governed by their own terms and privacy
              policies, and we are not responsible for them.
            </p>
            <p>
              Our handling of personal information is described in our Privacy
              Policy, which forms part of these Terms.
            </p>
          </LegalSection>

          <LegalSection id="changes" title="8. Changes to the Service and Terms">
            <p>
              We may add, change, or discontinue parts of the Service at any
              time. We may also update these Terms from time to time. When we
              make material changes, we will update the “Last updated” date above
              and, where appropriate, provide additional notice. Your continued
              use of the Service after a change takes effect means you accept the
              updated Terms. If you do not agree, stop using the Service and, if
              you wish, close your account.
            </p>
          </LegalSection>

          <LegalSection id="suspension" title="9. Suspension and termination">
            <p>
              You may stop using the Service and close your account at any time.
              We may suspend or terminate your access, with or without notice, if
              we reasonably believe you have violated these Terms or applicable
              law, if required by law or a third-party request, to protect the
              Service or other users, or for prolonged inactivity.
            </p>
            <p>
              On termination, your right to use the Service ends immediately.
              Sections that by their nature should survive (including credits
              already consumed, intellectual property, disclaimers, limitation of
              liability, indemnification, and governing law) will survive
              termination.
            </p>
          </LegalSection>

          <LegalSection id="disclaimers" title="10. Disclaimers">
            <p>
              The Service is provided “as is” and “as available”, without
              warranties of any kind, whether express, implied, statutory, or
              otherwise, including any implied warranties of merchantability,
              fitness for a particular purpose, title, non-infringement, or that
              the Service will be uninterrupted, error-free, secure, or that any
              Outputs will be accurate or reliable. Some jurisdictions do not
              allow the exclusion of certain warranties; in that case, such
              warranties are limited to the minimum extent permitted by law.
            </p>
          </LegalSection>

          <LegalSection id="liability" title="11. Limitation of liability">
            <p>
              To the maximum extent permitted by law, in no event will PT DXP
              Global Inovasi or its affiliates, officers, employees, agents,
              suppliers, or licensors be liable for any indirect, incidental,
              special, consequential, punitive, or exemplary damages, or for loss
              of profits, revenue, data, goodwill, or business opportunities,
              arising out of or relating to the Service or these Terms, even if
              we have been advised of the possibility of such damages.
            </p>
            <p>
              Our aggregate liability for any claims arising out of or relating
              to the Service or these Terms is limited to the greater of (a) the
              amount you paid us for credits in the twelve (12) months
              immediately preceding the event giving rise to the claim, or (b)
              one hundred US dollars (USD $100).
            </p>
            <p>
              Nothing in these Terms limits liability that cannot be limited under
              applicable law, including liability for fraud, gross negligence,
              willful misconduct, or death or personal injury caused by
              negligence.
            </p>
          </LegalSection>

          <LegalSection id="indemnity" title="12. Indemnification">
            <p>
              You agree to defend, indemnify, and hold harmless PT DXP Global
              Inovasi and its affiliates from and against any claims,
              liabilities, damages, losses, and expenses (including reasonable
              legal fees) arising out of or related to (a) your use of the
              Service, (b) your Inputs or Outputs, (c) your violation of these
              Terms or applicable law, or (d) your violation of any third-party
              right. We may assume the exclusive defense and control of any
              matter subject to indemnification, in which case you agree to
              cooperate with our defense.
            </p>
          </LegalSection>

          <LegalSection id="disputes" title="13. Governing law and disputes">
            <p>
              These Terms and any dispute arising out of or relating to them or
              the Service will be governed by the laws of the Republic of
              Indonesia, without regard to its conflict-of-laws rules. Where you
              reside in a jurisdiction whose consumer protection laws grant you
              mandatory rights, nothing in these Terms is intended to override
              those rights.
            </p>
            <p>
              The parties will first try to resolve any dispute informally by
              contacting us at{" "}
              <a href="mailto:hello@fromdxp.com" className="text-flame">
                hello@fromdxp.com
              </a>
              . If a dispute is not resolved within thirty (30) days, either
              party may pursue formal proceedings in a court of competent
              jurisdiction in Jakarta, Indonesia, except where applicable law
              requires a different forum or procedure.
            </p>
          </LegalSection>

          <LegalSection id="misc" title="14. Miscellaneous">
            <ul className="list-disc space-y-2 pl-5">
              <li>
                <strong>Entire agreement.</strong> These Terms, together with the
                Privacy Policy and any other policies referenced here, are the
                entire agreement between you and us regarding the Service.
              </li>
              <li>
                <strong>Severability.</strong> If any provision is held invalid
                or unenforceable, the remaining provisions remain in effect.
              </li>
              <li>
                <strong>No waiver.</strong> Our failure to enforce a provision is
                not a waiver of our right to do so later.
              </li>
              <li>
                <strong>Assignment.</strong> You may not assign these Terms
                without our prior written consent. We may assign them in
                connection with a merger, acquisition, or sale of assets.
              </li>
              <li>
                <strong>Notices.</strong> We may give notices through the Service
                or by email to the address associated with your account.
              </li>
              <li>
                <strong>Force majeure.</strong> We are not liable for delay or
                failure to perform due to events beyond our reasonable control.
              </li>
              <li>
                <strong>No agency.</strong> Nothing in these Terms creates any
                agency, partnership, joint venture, or employment relationship.
              </li>
            </ul>
          </LegalSection>

          <LegalSection id="contact" title="15. Contact us">
            <p>
              If you have questions about these Terms, contact us at{" "}
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
