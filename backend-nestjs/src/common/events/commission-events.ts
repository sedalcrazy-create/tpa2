/**
 * Event Contract: Commission Module ↔ TPA Core
 * Version: 1.0.0
 *
 * This file defines the event schemas for communication between
 * the Medical Commission module (NestJS) and TPA Core engine (Go).
 *
 * Events are:
 * - Versioned (schema version in each event)
 * - Immutable (once published, never modified)
 * - Auditable (event_id, timestamp, correlation)
 * - Transport-agnostic (JSON)
 */

export const EVENT_VERSION = '1.0.0';

// Event Types
export enum EventType {
  // Commission case events
  COMMISSION_CASE_CREATED = 'commission.case.created',
  COMMISSION_CASE_ASSIGNED = 'commission.case.assigned',
  COMMISSION_CASE_REVIEWED = 'commission.case.reviewed',
  COMMISSION_VERDICT_ISSUED = 'commission.verdict.issued',
  COMMISSION_CASE_CLOSED = 'commission.case.closed',
  COMMISSION_CASE_APPEALED = 'commission.case.appealed',

  // Social work events
  SOCIALWORK_CASE_CREATED = 'socialwork.case.created',
  SOCIALWORK_ASSESSMENT_DONE = 'socialwork.assessment.done',
  SOCIALWORK_REFERRAL_ISSUED = 'socialwork.referral.issued',
}

// Base event interface
export interface BaseEvent {
  event_id: string; // UUID
  event_type: EventType;
  version: string;
  timestamp: string; // ISO 8601
  source: string; // 'commission-api' | 'tpa-api'

  // Correlation
  correlation_id?: string;
  causation_id?: string;

  // Tenant context
  tenant_id: number;
}

// Insured person reference
export interface InsuredPersonRef {
  id: string;
  national_id: string;
  personnel_code: string;
  full_name: string;
  relation: 'main' | 'spouse' | 'child' | 'parent';
}

// Verdict details
export interface VerdictDetails {
  verdict_id: string;
  verdict_type: string;
  verdict_code: string;
  verdict_text: string;
  disability_rate?: number; // 0-100
  effective_from?: string;
  effective_to?: string;
  is_permanent: boolean;
  needs_review: boolean;
  review_date?: string;
  approved_by: string;
  approved_at: string;
  commission_level: 'provincial' | 'central';
}

// Financial impact
export interface FinancialImpact {
  coverage_type: 'full' | 'partial' | 'none';
  coverage_percent: number;
  coverage_limit_delta: number;

  monthly_allowance?: number;
  lump_sum_payment?: number;

  service_restrictions?: string[];
  provider_restrictions?: string[];

  valid_from: string;
  valid_to?: string;
}

// Document reference
export interface DocumentRef {
  document_id: string;
  document_type: 'verdict_letter' | 'medical_report' | 'attachment';
  file_name: string;
  file_url: string;
  checksum: string;
}

// Commission Verdict Event
export interface CommissionVerdictEvent extends BaseEvent {
  event_type: EventType.COMMISSION_VERDICT_ISSUED;

  case_id: string;
  case_number: string;

  insured_person: InsuredPersonRef;
  verdict: VerdictDetails;
  financial_impact?: FinancialImpact;
  documents?: DocumentRef[];
}

// Referral details for social work
export interface ReferralDetails {
  referral_id: string;
  referral_type: string;
  referral_reason: string;
  referred_to: string;
  priority: 'normal' | 'urgent' | 'critical';
  due_date?: string;
  amount?: number;
  issued_by: string;
  issued_at: string;
  notes?: string;
}

// Social Work Referral Event
export interface SocialWorkReferralEvent extends BaseEvent {
  event_type: EventType.SOCIALWORK_REFERRAL_ISSUED;

  case_id: string;
  case_type: string;

  insured_person: InsuredPersonRef;
  referral: ReferralDetails;
}

// Event envelope for transport
export interface EventEnvelope<T extends BaseEvent> {
  event: T;
  schema: string;
  encoded: boolean;
}

// Helper to create events
export function createEvent<T extends BaseEvent>(
  type: EventType,
  tenantId: number,
  payload: Omit<T, keyof BaseEvent>,
): T {
  const base: BaseEvent = {
    event_id: crypto.randomUUID(),
    event_type: type,
    version: EVENT_VERSION,
    timestamp: new Date().toISOString(),
    source: 'commission-api',
    tenant_id: tenantId,
  };

  return { ...base, ...payload } as T;
}

// Validation
export function validateEvent(event: BaseEvent): string[] {
  const errors: string[] = [];

  if (!event.event_id) errors.push('event_id is required');
  if (!event.tenant_id) errors.push('tenant_id is required');
  if (!event.event_type) errors.push('event_type is required');
  if (!event.timestamp) errors.push('timestamp is required');

  return errors;
}
