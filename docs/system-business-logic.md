# Lost & Found Ethiopia — Business Logic → Entity Derivation

Purpose: define the business rules first, then derive entities/relationships from them. This keeps the schema grounded in real behavior instead of guessed tables.

Scope note: rules marked **[MVP]** are needed for the Addis Ababa Phase 1 launch. Rules marked **[Phase 2+]** can be designed for but not necessarily built yet — good schemas anticipate them without over-engineering now.

---

## 1. User & Identity Logic

- A user can act as a **reporter** (of a lost or found item), a **claimant**, or both — same person, different roles per report. **[MVP]**
- A user must have a **verified fayda id** (via SMS OTP) before they can post a report or file a claim — fayda is the primary trust anchor in Ethiopia (more reliable than others) if does not have fayda id use the below which is phone. **[MVP]**
- A user must have a **verified phone number** (via SMS OTP) before they can post a report or file a claim — phone is the primary trust anchor in Ethiopia (more reliable than email). **[MVP]**
- A user has a **preferred language** (Amharic/English) but later we make to support (Oromo/Tigrinya/) in the next versions that drives notification content. as well the over all app content **[MVP]**
- A user can optionally upgrade to an **institutional account** (hotel, university, transport company) with different permissions and a public verified badge. **[Phase 2]**
- A user has a **trust score** that increases with successful resolutions and decreases with flagged/fraudulent claims — this score should influence how much friction they face in future claims (e.g., low-trust users get stricter verification). **[Phase 2, but reserve the field now]**
- Users can be **suspended/banned** by admins — must retain their historical reports for audit even if banned (soft-delete, not hard-delete).

→ **Entities implied:** `User`, `UserRole` (or role enum on User), `Institution` (linked 1:1 to a User account that manages it), `TrustEvent` (log of what changed a user's trust score).

---

## 2. Reporting Logic (Lost & Found)

- A report is always one of exactly two types: `LOST` or `FOUND` — never both, never neither. **[MVP]**
- A report must belong to exactly one **category** (ID document, phone, wallet, pet, vehicle, other), and category determines which extra fields are required (e.g., IMEI for phones, microchip ID for pets, plate number for vehicles). **[MVP]**
- A report requires: description, at least one photo (optional for IDs where photographing sensitive info is discouraged), approximate location (lat/lng or named landmark), and a date/time window (not always an exact moment — "sometime yesterday afternoon" is realistic). **[MVP]**
- A report has a **status lifecycle**: `OPEN → MATCHED → CLAIM_PENDING → RESOLVED` or `OPEN → EXPIRED/CLOSED`. Status transitions are business-critical and must be logged. **[MVP]**

- A reporter can **edit or withdraw** their own report only while it's still `OPEN` (once matched/claim-pending, changes need admin review to prevent tampering with a live dispute). **[MVP]**
- For **government ID reports specifically**: the ID number is never stored in full — only a masked/hashed version, used for exact-match verification without ever displaying it back. **[MVP — trust-critical]**

→ **Entities implied:** `Report` (type, status, description, dates), `Category` (with a schema of which extra fields it requires), `ReportAttribute` (flexible key-value store per report for category-specific fields, or category-specific subtables), `ReportStatusHistory`.

---

## 3. Category & Attribute Logic

- Categories are **not hardcoded in application logic** — they should be data-driven (a `Category` table with a defined set of expected fields) so admins can add new categories (e.g., "luggage," "laptops") without a code deploy. **[MVP — design for this even if only 4-5 categories launch]**
- Each category defines: required fields, whether photos are mandatory, whether the item can have a unique identifier (IMEI, plate number, microchip, ID number) used for exact matching. **[MVP]**
- Categories may have **priority/urgency weighting** — e.g., ID and phone reports get faster SMS alert cycles than "other" category, because urgency is genuinely higher. **[Phase 2]**

→ **Entities implied:** `Category`, `CategoryFieldDefinition` (defines the dynamic schema per category).

---

## 4. Matching Logic

- A **Match** is a system-generated (or admin-generated) link between one `LOST` report and one `FOUND` report, with a **confidence score** based on category match, location proximity, date-range overlap, and description/keyword similarity. **[MVP — start rule-based, not ML]**
- A single lost report can have **multiple candidate matches**; only one can become `CONFIRMED` at a time. Once confirmed, other candidate matches for the same report should be auto-closed. **[MVP]**
- Matching runs **automatically** on new report creation (compare new report against the opposite-type pool) and can also be **triggered manually** by a user browsing found listings and self-identifying a match. **[MVP]**
- A match does **not** immediately reveal contact info to either party — it only triggers the **claim/verification flow** (see next section). Revealing raw contact info without verification is a core safety violation of this business's whole value proposition. **[MVP — non-negotiable rule]**
- For exact-identifier categories (IMEI, plate number, ID number, microchip), an exact match on the hashed identifier should short-circuit to very high confidence — this is your most reliable signal type. **[MVP]**

→ **Entities implied:** `Match` (lost_report_id, found_report_id, confidence_score, status, created_at), `MatchStatusHistory`.

---

## 5. Claim & Verification Logic (the trust core of the business)

- Only the person who filed the **Lost** report (or, in institutional flows, a person claiming to be the owner via other means) can initiate a **Claim** against a Match. **[MVP]**
- A claim requires answering **challenge questions** defined by the finder/reporter at report time (e.g., "what's the last digit of the ID," "what color is the phone case," "pet's name") — these answers are compared, not publicly stored as the "correct answer" anywhere the claimant can see. **[MVP]**
- A claim has a **pass/fail outcome** per attempt, and the system should **rate-limit claim attempts** (e.g., max 3 tries per match, then requires admin manual review) to prevent brute-forcing. **[MVP — critical anti-fraud rule]**
- Only after a claim is **verified** (either by passing challenge questions or by admin manual approval) does the system reveal **masked contact details** and unlock the in-app chat for handover coordination. **[MVP]**
- For **government ID categories**, a successful challenge-question pass is not sufficient alone — the business logic should flag these for **mandatory admin review** or in-person Kebele/police verification before full contact release, given the legal sensitivity. **[MVP — policy rule, not just code]**
- Every claim attempt (pass or fail) must be **logged immutably** for dispute resolution and fraud pattern detection. **[MVP]**
- A finder can **decline a claim** if they believe it's not genuine, which reopens the match for other candidates or escalates to admin. **[MVP]**

→ **Entities implied:** `Claim` (match_id, claimant_id, status, attempt_count), `ChallengeQuestion` (per report, defined by reporter), `ChallengeAnswer` (submitted per claim attempt), `ClaimAttemptLog`.

---

## 6. Communication Logic

- Chat is only unlocked **after a claim is verified** — never before. Messages should reference the match, not exist as freestanding DMs. **[MVP]**
- Phone numbers stay **masked** inside the chat (use of a proxy/relay number or simply a rule that raw numbers typed in chat get flagged) until both parties explicitly agree to share direct contact for the physical handover. **[Phase 2 for masking infra; MVP can start simpler — just strongly encourage in-app coordination]**
- All messages tied to a match should be **retained** even after resolution, for dispute history. **[MVP]**

→ **Entities implied:** `Message` (match_id, sender_id, content, timestamp), `Conversation` (optional wrapper if you want to group messages beyond 1:1 match scope later).

---

## 7. Resolution & Closure Logic

- A match moves to `RESOLVED` only when **both parties confirm** the handover happened (a simple "mark as returned" action from both sides is a strong, simple rule to start with). **[MVP]**
- A resolved match should prompt an optional **rating/feedback** from both sides — this feeds the trust score system later and gives you your key success metric (resolution rate). **[MVP: capture the confirmation. Feedback/ratings can be Phase 2.]**
- If only one party confirms after a set time window, the system should **nudge the other party** via notification before escalating to "disputed" status for admin review. **[Phase 2]**

→ **Entities implied:** `ResolutionConfirmation` (match_id, user_id, confirmed_at), `Feedback` (Phase 2).

---

## 8. Institutional Logic

- Institutions (hotels, universities, airport, ride-hailing) can **bulk-create found reports** and manage their own internal "desk" of held items. **[Phase 2]**
- Institutional accounts have a **verified badge**, driven by an admin-approved verification process (business registration check), not self-declared. **[Phase 2]**
- Institutions may have **multiple staff logins** under one institutional account — meaning permissions need to be scoped at the institution level, not just the individual user level. **[Phase 2 — but worth anticipating in the User/Institution relationship now]**

→ **Entities implied:** `Institution`, `InstitutionMember` (join table: user_id, institution_id, role_within_institution).

---

## 9. Notification Logic

- Every meaningful state transition (new match found, claim submitted, claim verified, chat message received, resolution confirmed, report about to expire) triggers a notification. **[MVP]**
- Notification **channel** depends on user preference and what's available — SMS is the reliable fallback; push/email are enhancements, not the baseline, given device/connectivity variance. **[MVP]**
- Users should be able to control notification frequency/channel per category of alert (e.g., always SMS for ID matches, push-only for "other" category) — but a simple global on/off is fine for MVP. **[MVP: simple. Phase 2: granular.]**

→ **Entities implied:** `Notification` (user_id, type, channel, payload, status: sent/delivered/failed, related_entity_id).

---

## 10. Admin & Moderation Logic

- Admins can **view, filter, and act on** any report, match, or claim — with every action logged (who did what, when) for accountability. **[MVP]**
- Admins can **force-verify** or **force-reject** a claim (overriding the automated challenge-question flow) — needed for edge cases like ID document claims requiring human judgment. **[MVP]**
- Admins can **flag/ban users** and must be able to see a user's full history (reports, claims, past flags) in one view before acting. **[MVP]**
- A separate, smaller "moderator" role should be possible without full admin/system access (principle of least privilege) — even if only one person uses the platform today, this prevents painful re-architecture later. **[Design now, enforce loosely at MVP]**

→ **Entities implied:** `AdminActionLog` (admin_id, action_type, target_entity, target_id, timestamp, notes).

---

## 11. Trust & Fraud Prevention Logic (cross-cutting)

- Any entity that can be **abused** (reports, claims, messages) needs a `flagged` boolean and a relation to a `Report/Flag` record explaining why and by whom. **[MVP]**
- Rate limits apply at multiple levels: reports per user per day, claim attempts per match, messages per minute — these are business rules, not just infra config, because they protect the platform's core trust promise. **[MVP]**
- Every sensitive action (claim attempt, contact reveal, admin override) should be **immutably logged** — treat this like a financial ledger, not just a debug log. **[MVP]**

→ **Entities implied:** `Flag` (target_type, target_id, reason, raised_by, resolved_by), `RateLimitEvent` (or handled at infra/cache layer, but the rule itself is a business decision worth documenting).

---

## Consolidated Candidate Entity List (preview — ERD comes next)

| Entity | Core Purpose |
|---|---|
| `User` | Any person using the platform |
| `Institution` | Business/org accounts |
| `InstitutionMember` | Join: users ↔ institutions |
| `Category` | Item type definitions |
| `CategoryFieldDefinition` | Dynamic required fields per category |
| `Report` | A lost or found item report |
| `ReportAttribute` | Category-specific field values per report |
| `ReportStatusHistory` | Audit trail of report status changes |
| `ChallengeQuestion` | Verification questions set by reporter |
| `Match` | Link between a lost and found report |
| `MatchStatusHistory` | Audit trail of match status changes |
| `Claim` | A claim attempt against a match |
| `ChallengeAnswer` | Submitted answers per claim attempt |
| `ClaimAttemptLog` | Immutable log of claim attempts |
| `Message` | Chat tied to a verified match |
| `ResolutionConfirmation` | Both-side confirmation of handover |
| `Feedback` | Post-resolution rating (Phase 2) |
| `Notification` | Outbound alerts to users |
| `Flag` | Fraud/abuse reports on any entity |
| `AdminActionLog` | Accountability trail for admin actions |
| `TrustEvent` | What changed a user's trust score (Phase 2) |

---

## Next Step

This list gives us clean entities because every one of them traces back to a real business rule rather than a guess. Next: turn this into a proper **ER diagram + full schema** (tables, columns, types, foreign keys, indexes) — including the tricky parts like how `ReportAttribute` handles dynamic category fields (e.g., EAV pattern vs. JSONB column vs. per-category subtables, which is a real design decision worth walking through).




# Lost & Found Ethiopia — Formal ER Design (Entities → Relationships → Attributes)

Methodology: **Conceptual Design → Logical Design → Physical Design**, the standard professional database design process.
- Step 1: identify entities (conceptual)
- Step 2: define relationships & cardinality between them (conceptual)
- Step 3: define attributes per entity, with keys and types (logical, ready for physical implementation)

---

# STEP 1 — Entity Identification

Every entity below traces back to a rule in our business logic document — nothing here is guessed.

| # | Entity | Type | One-line definition |
|---|---|---|---|
| 1 | **User** | Core | Any person using the platform (reporter, claimant, or admin) |
| 2 | **Institution** | Core | A business/organization account (hotel, university, transport co.) |
| 3 | **InstitutionMember** | Associative | Links Users to Institutions they belong to, with a role |
| 4 | **Category** | Reference | An item type (Government ID, Phone, Pet, Vehicle, etc.) |
| 5 | **CategoryFieldDefinition** | Reference | Defines what dynamic fields a Category requires |
| 6 | **Report** | Core | A lost or found item report |
| 7 | **ReportStatusHistory** | Audit | Log of a Report's status changes over time |
| 8 | **ChallengeQuestion** | Core | A verification question set by a reporter |
| 9 | **Match** | Core | A system/admin-suggested link between a lost and a found Report |
| 10 | **MatchStatusHistory** | Audit | Log of a Match's status changes over time |
| 11 | **Claim** | Core | An attempt by a claimant to verify ownership against a Match |
| 12 | **ClaimAttemptLog** | Audit | Immutable log of each challenge-question attempt |
| 13 | **Message** | Core | A chat message tied to a verified Match |
| 14 | **ResolutionConfirmation** | Core | One party's confirmation that a handover happened |
| 15 | **Feedback** | Core (Phase 2) | Post-resolution rating between two Users |
| 16 | **Notification** | Support | A record of an alert sent to a User |
| 17 | **Flag** | Support | A fraud/abuse report raised against any entity |
| 18 | **AdminActionLog** | Audit | Accountability log of admin actions |
| 19 | **TrustEvent** | Support (Phase 2) | A record of what changed a User's trust score |

**Entity classification matters professionally** because it tells you how to treat each table operationally:
- **Core** entities are what the product is about — they change often, need strong validation.
- **Reference** entities are mostly static/admin-managed — cacheable, rarely written.
- **Associative** entities exist only to resolve many-to-many relationships.
- **Audit** entities are append-only, never updated or deleted — treated like a ledger.
- **Support** entities are operational plumbing (notifications, trust) — high write volume, low query complexity.

---

# STEP 2 — Relationships & Cardinality

Notation: `1:1` one-to-one, `1:N` one-to-many, `M:N` many-to-many (always resolved via an associative entity in the physical design).

| # | Relationship | Cardinality | Business justification |
|---|---|---|---|
| R1 | User — Institution | `M:N` (via InstitutionMember) | One user can belong to multiple institutions (rare but possible — e.g. a consultant); one institution has multiple staff |
| R2 | User — Report | `1:N` | A user can file many reports over time; each report has exactly one reporter |
| R3 | Institution — Report | `1:N` (optional) | An institution can log many found items; a report may or may not be tied to an institution |
| R4 | Category — Report | `1:N` | A category can apply to many reports; each report belongs to exactly one category |
| R5 | Category — CategoryFieldDefinition | `1:N` | A category defines multiple dynamic fields |
| R6 | Report — ReportStatusHistory | `1:N` | A report accumulates a history of status changes over its lifecycle |
| R7 | Report — ChallengeQuestion | `1:N` | A reporter can set multiple verification questions per report |
| R8 | Report — Match | `1:N` (as lost) and `1:N` (as found) | A single lost report can have multiple candidate found matches, and vice versa, until one is confirmed |
| R9 | Match — MatchStatusHistory | `1:N` | Same audit pattern as reports |
| R10 | Match — Claim | `1:N` | A match can receive multiple claim attempts (e.g., rejected claim, then a legitimate one) |
| R11 | User — Claim | `1:N` | A user can be the claimant on multiple claims (across different matches) |
| R12 | Claim — ChallengeQuestion (via ClaimAttemptLog) | `M:N` | A claim attempts to answer multiple questions; each question can be attempted across multiple claims (if a report gets multiple claim attempts) |
| R13 | Match — Message | `1:N` | A match has a chat thread of many messages |
| R14 | User — Message | `1:N` | A user can send many messages (as sender) |
| R15 | Match — ResolutionConfirmation | `1:N` (max 2 in practice — one per party) | Both parties confirm independently |
| R16 | Match — Feedback | `1:N` | Both parties can each leave feedback about the other |
| R17 | User — Notification | `1:N` | A user receives many notifications over time |
| R18 | User — Flag | `1:N` (as raiser) | Any entity can be flagged; flags reference their target generically (polymorphic association) |
| R19 | User — AdminActionLog | `1:N` (as admin) | An admin performs many logged actions |
| R20 | User — TrustEvent | `1:N` | A user accumulates many trust-affecting events over time |

**Why this step matters before writing attributes:** cardinality decisions directly determine where foreign keys live. A `1:N` means the foreign key sits on the "many" side (e.g., `reports.reporter_id` references `users.id` — not the other way around). Getting this order right *before* touching column types avoids the classic beginner mistake of foreign keys pointing the wrong direction or unnecessary junction tables for what's actually a simple `1:N`.

---

# STEP 3 — Attributes per Entity (Data Dictionary)

Now that entities and relationships are fixed, each entity gets its attributes, with **PK** (Primary Key), **FK** (Foreign Key), **UK** (Unique Key) marked explicitly.

### 3.1 User
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | Unique identifier |
| phone_number | VARCHAR(20) | UK | No | Primary trust anchor / login identifier |
| phone_verified_at | TIMESTAMPTZ | — | Yes | Null until OTP-verified |
| email | VARCHAR(255) | UK | Yes | Optional secondary contact |
| full_name | VARCHAR(255) | — | No | |
| password_hash | TEXT | — | Yes | Null if phone-OTP-only auth |
| preferred_language | VARCHAR(5) | — | No | am / om / ti / en |
| trust_score | INTEGER | — | No | Default 100, adjusted by TrustEvent |
| status | ENUM | — | No | active / suspended / banned |
| created_at / updated_at | TIMESTAMPTZ | — | No | |

### 3.2 Institution
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| name | VARCHAR(255) | — | No | |
| type | ENUM | — | No | hotel / university / transport / bank / other |
| business_reg_number | VARCHAR(100) | — | Yes | For verification workflow |
| verified | BOOLEAN | — | No | Admin-approved only |
| verified_by | UUID | FK → User.id | Yes | Which admin verified it |
| location | GEOGRAPHY(Point) | — | Yes | PostGIS point |
| created_at | TIMESTAMPTZ | — | No | |

### 3.3 InstitutionMember
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| institution_id | UUID | FK → Institution.id | No | |
| user_id | UUID | FK → User.id | No | |
| role | ENUM | — | No | owner / staff |
| created_at | TIMESTAMPTZ | — | No | |
| *(institution_id, user_id)* | — | UK (composite) | — | Prevents duplicate memberships |

### 3.4 Category
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| slug | VARCHAR(50) | UK | No | e.g. 'government_id' |
| name_am / name_en / name_om / name_ti | VARCHAR(100) | — | Mixed | Localized display names |
| requires_unique_identifier | BOOLEAN | — | No | Drives exact-match logic |
| requires_photo | BOOLEAN | — | No | |
| priority_weight | SMALLINT | — | No | Urgency ranking |
| is_active | BOOLEAN | — | No | |
| created_at | TIMESTAMPTZ | — | No | |

### 3.5 CategoryFieldDefinition
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| category_id | UUID | FK → Category.id | No | |
| field_key | VARCHAR(50) | — | No | e.g. 'imei' |
| field_label_am / field_label_en | VARCHAR(100) | — | No | |
| field_type | VARCHAR(20) | — | No | text / number / date |
| is_required | BOOLEAN | — | No | |
| is_unique_identifier | BOOLEAN | — | No | Marks the field used for exact-match |
| *(category_id, field_key)* | — | UK (composite) | — | |

### 3.6 Report
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| reporter_id | UUID | FK → User.id | No | |
| institution_id | UUID | FK → Institution.id | Yes | Optional |
| type | ENUM | — | No | lost / found |
| category_id | UUID | FK → Category.id | No | |
| status | ENUM | — | No | open / matched / claim_pending / resolved / expired / closed |
| title | VARCHAR(255) | — | No | |
| description | TEXT | — | No | |
| attributes | JSONB | — | No | Category-specific fields |
| identifier_hash | TEXT | — | Yes | Hashed unique identifier for exact matching |
| photo_urls | TEXT[] | — | Yes | Cloudinary URLs |
| location | GEOGRAPHY(Point) | — | Yes | |
| location_label | VARCHAR(255) | — | Yes | Human-readable location |
| event_date_start / event_date_end | DATE | — | No | Supports date ranges |
| expires_at | TIMESTAMPTZ | — | No | Default now() + 90 days |
| created_at / updated_at | TIMESTAMPTZ | — | No | |

### 3.7 ReportStatusHistory
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| report_id | UUID | FK → Report.id | No | |
| from_status / to_status | ENUM | — | Mixed | from_status null on creation |
| changed_by | UUID | FK → User.id | Yes | Null if system-triggered |
| reason | TEXT | — | Yes | |
| created_at | TIMESTAMPTZ | — | No | |

### 3.8 ChallengeQuestion
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| report_id | UUID | FK → Report.id | No | |
| question | TEXT | — | No | |
| answer_hash | TEXT | — | No | Never plaintext |
| created_at | TIMESTAMPTZ | — | No | |

### 3.9 Match
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| lost_report_id | UUID | FK → Report.id | No | |
| found_report_id | UUID | FK → Report.id | No | |
| confidence_score | NUMERIC(5,2) | — | No | 0.00–100.00 |
| match_reason | JSONB | — | No | Explains scoring (distance, exact ID hit, etc.) |
| status | ENUM | — | No | suggested / confirmed / rejected / expired |
| created_at / updated_at | TIMESTAMPTZ | — | No | |
| *(lost_report_id, found_report_id)* | — | UK (composite) | — | Prevents duplicate match rows |

### 3.10 MatchStatusHistory
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| match_id | UUID | FK → Match.id | No | |
| from_status / to_status | ENUM | — | Mixed | |
| changed_by | UUID | FK → User.id | Yes | |
| created_at | TIMESTAMPTZ | — | No | |

### 3.11 Claim
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| match_id | UUID | FK → Match.id | No | |
| claimant_id | UUID | FK → User.id | No | |
| status | ENUM | — | No | pending / verified / rejected / admin_review |
| attempt_count | SMALLINT | — | No | Default 0, incremented per attempt |
| requires_admin_review | BOOLEAN | — | No | True by policy for government_id category |
| reviewed_by | UUID | FK → User.id | Yes | Admin who reviewed, if applicable |
| reviewed_at | TIMESTAMPTZ | — | Yes | |
| created_at / updated_at | TIMESTAMPTZ | — | No | |

### 3.12 ClaimAttemptLog
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| claim_id | UUID | FK → Claim.id | No | |
| question_id | UUID | FK → ChallengeQuestion.id | No | |
| submitted_answer_hash | TEXT | — | No | |
| passed | BOOLEAN | — | No | |
| attempted_at | TIMESTAMPTZ | — | No | |
| ip_address | INET | — | Yes | Fraud pattern detection |

### 3.13 Message
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| match_id | UUID | FK → Match.id | No | |
| sender_id | UUID | FK → User.id | No | |
| content | TEXT | — | No | |
| flagged | BOOLEAN | — | No | e.g. auto-flagged if raw phone number detected |
| created_at | TIMESTAMPTZ | — | No | |

### 3.14 ResolutionConfirmation
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| match_id | UUID | FK → Match.id | No | |
| confirmed_by | UUID | FK → User.id | No | |
| confirmed_at | TIMESTAMPTZ | — | No | |
| *(match_id, confirmed_by)* | — | UK (composite) | — | One confirmation per user per match |

### 3.15 Feedback (Phase 2)
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| match_id | UUID | FK → Match.id | No | |
| from_user_id | UUID | FK → User.id | No | |
| about_user_id | UUID | FK → User.id | No | |
| rating | SMALLINT | — | No | 1–5, CHECK constraint |
| comment | TEXT | — | Yes | |
| created_at | TIMESTAMPTZ | — | No | |

### 3.16 Notification
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| user_id | UUID | FK → User.id | No | |
| channel | ENUM | — | No | sms / push / email |
| type | VARCHAR(50) | — | No | e.g. 'match_found' |
| payload | JSONB | — | No | |
| related_entity_type | VARCHAR(30) | — | Yes | 'report' / 'match' / 'claim' |
| related_entity_id | UUID | — | Yes | Polymorphic reference |
| status | ENUM | — | No | pending / sent / delivered / failed |
| sent_at | TIMESTAMPTZ | — | Yes | |
| created_at | TIMESTAMPTZ | — | No | |

### 3.17 Flag
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| target_type | VARCHAR(30) | — | No | 'report' / 'claim' / 'message' / 'user' |
| target_id | UUID | — | No | Polymorphic reference (enforced in app layer) |
| raised_by | UUID | FK → User.id | Yes | Null if system-auto-flagged |
| reason | TEXT | — | No | |
| resolved | BOOLEAN | — | No | |
| resolved_by | UUID | FK → User.id | Yes | |
| resolved_at | TIMESTAMPTZ | — | Yes | |
| created_at | TIMESTAMPTZ | — | No | |

### 3.18 AdminActionLog
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| admin_id | UUID | FK → User.id | No | |
| action_type | VARCHAR(50) | — | No | e.g. 'force_verify_claim' |
| target_type | VARCHAR(30) | — | No | |
| target_id | UUID | — | No | |
| notes | TEXT | — | Yes | |
| created_at | TIMESTAMPTZ | — | No | |

### 3.19 TrustEvent (Phase 2)
| Attribute | Type | Key | Nullable | Description |
|---|---|---|---|---|
| id | UUID | PK | No | |
| user_id | UUID | FK → User.id | No | |
| delta | SMALLINT | — | No | +5 / -20 etc. |
| reason | VARCHAR(100) | — | No | e.g. 'resolution_confirmed' |
| related_entity_id | UUID | — | Yes | |
| created_at | TIMESTAMPTZ | — | No | |

---

# Normalization Check (professional sanity check)

- **1NF:** every attribute is atomic — the one deliberate exception is `Report.attributes` (JSONB) and `Report.photo_urls` (array), which are intentional denormalizations for flexibility, not oversights. This is a documented, deliberate trade-off, not a mistake.
- **2NF:** every non-key attribute depends on the *whole* primary key — satisfied, since every table uses a single-column UUID surrogate key (no partial-key dependency issues possible).
- **3NF:** no non-key attribute depends on another non-key attribute — satisfied. E.g., `Category` names aren't duplicated onto `Report`; `Report` only stores `category_id` and looks up names via join.

This schema is in 3NF except for the two explicitly justified JSONB/array exceptions — which is exactly what a professional reviewer would want to see: normalized by default, denormalized only where there's a clear, stated reason.

---
chi — solid pick. It's idiomatic (stays close to net/http instead of wrapping everything in framework magic), has clean middleware chaining, and pairs very well with gorilla/websocket for the chat feature we flagged as must-have. It's also easier to reason about long-term than Gin/Fiber if the team grows, since it doesn't lock you into a non-standard request/response type.

A few things worth setting up from day one with chi, since they map straight onto the business rules we defined:

Middleware chain: auth (JWT) → rate limiter (Redis-backed, for the claim-attempt and report-submission limits) → logger → recoverer. chi's middleware composition makes this clean.
Route grouping by resource (/reports, /matches, /claims, /admin/*) with chi.Router.Route() — keeps the admin-only routes cleanly separated with their own auth middleware.
For WebSocket chat, mount it as its own route (/ws/chat/{matchID}) with the same auth middleware applied, since chat should only be reachable after a claim is verified — enforce that at the handshake, not just in the frontend.

Cloudinary — good fit, especially at MVP stage.

Built-in image transformation/compression on upload is genuinely useful here — you want thumbnails for the browse/search list and full-res only on the detail view, and Cloudinary does that via URL params without you writing a resize pipeline.
It also has moderation/auto-tagging features (AI content moderation) which could be a nice extra layer against people uploading irrelevant or inappropriate images.
One thing to plan for: Cloudinary's free/lower tiers have bandwidth and storage caps — fine for MVP and early growth, but worth knowing the pricing cliff before you scale, so it doesn't surprise you later. If you ever outgrow it, migrating to Cloudflare R2 + your own transformation logic is straightforward since URLs are abstracted behind your API anyway (don't let the frontend hardcode Cloudinary URLs directly — proxy through your own image field so swapping providers later doesn't touch your DB or frontend code).

Your stack is now fully locked in:
Go (chi) + PostgreSQL (PostGIS, JSONB) + Redis (Asynq + rate limiting + WS pub/sub) + Cloudinary + Next.js + Flutter

That's a genuinely professional, appropriately-scoped stack for this