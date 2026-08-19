import { useEffect, useState } from "react";

import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  MenuItem,
  Stack,
  TextField,
} from "@mui/material";

import {
  createLoan,
  updateLoan,
} from "../../services/loanService";

import type {
  CreateLoanRequest,
  Loan,
  UpdateLoanRequest,
} from "../../types/loan";

interface Props {
  open: boolean;
  loan?: Loan | null;
  onClose: () => void;
  onSuccess: () => void;
}

interface FormData {
  client_id: string;
  pn_number: string;
  loan_type: string;
  principal_amount: string;
  interest_rate: string;
  loan_interest: string;
  pn_value: string;
  loan_term: string;
  amortization_amount: string;
  disbursement_date: string;
  maturity_date: string;
  frequency: string;
  status: string;
}

const emptyForm: FormData = {
  client_id: "",
  pn_number: "",
  loan_type: "",
  principal_amount: "",
  interest_rate: "",
  loan_interest: "",
  pn_value: "",
  loan_term: "",
  amortization_amount: "",
  disbursement_date: "",
  maturity_date: "",
  frequency: "MONTHLY",
  status: "ACTIVE",
};

export default function LoanFormDialog({
  open,
  loan,
  onClose,
  onSuccess,
}: Props) {
  const [form, setForm] =
    useState<FormData>(emptyForm);

  const [saving, setSaving] =
    useState(false);

  const [error, setError] =
    useState("");

  const editing = Boolean(loan);

  // ========================================
  // LOAD FORM DATA
  // ========================================

  useEffect(() => {
    if (!open) {
      return;
    }

    setError("");

    if (!loan) {
      setForm(emptyForm);
      return;
    }

    setForm({
      client_id: String(
        loan.client_id ?? "",
      ),

      pn_number:
        loan.pn_number ?? "",

      loan_type:
        loan.loan_type ?? "",

      principal_amount:
        String(
          loan.principal_amount ?? "",
        ),

      interest_rate:
        String(
          loan.interest_rate ?? "",
        ),

      loan_interest:
        String(
          loan.loan_interest ?? "",
        ),

      pn_value:
        String(
          loan.pn_value ?? "",
        ),

      loan_term:
        String(
          loan.loan_term ?? "",
        ),

      amortization_amount:
        String(
          loan.amortization_amount ?? "",
        ),

      disbursement_date:
        loan.disbursement_date
          ? loan.disbursement_date.slice(
              0,
              10,
            )
          : "",

      maturity_date:
        loan.maturity_date
          ? loan.maturity_date.slice(
              0,
              10,
            )
          : "",

      frequency:
        loan.frequency ||
        "MONTHLY",

      status:
        loan.status ||
        "ACTIVE",
    });
  }, [open, loan]);

  // ========================================
  // FIELD CHANGE
  // ========================================

  const change = (
    field: keyof FormData,
    value: string,
  ) => {
    setForm((current) => ({
      ...current,
      [field]: value,
    }));
  };

  // ========================================
  // SAVE
  // ========================================

  const handleSubmit = async () => {
    try {
      setSaving(true);
      setError("");

      if (!form.client_id) {
        throw new Error(
          "Client ID is required.",
        );
      }

      if (!form.pn_number.trim()) {
        throw new Error(
          "PN number is required.",
        );
      }

      if (
        Number(
          form.principal_amount,
        ) <= 0
      ) {
        throw new Error(
          "Principal amount must be greater than zero.",
        );
      }

      if (
        Number(form.loan_term) <= 0
      ) {
        throw new Error(
          "Loan term must be greater than zero.",
        );
      }

      const payload = {
        client_id: Number(
          form.client_id,
        ),

        pn_number:
          form.pn_number.trim(),

        loan_type:
          form.loan_type.trim() ||
          null,

        principal_amount: Number(
          form.principal_amount,
        ),

        interest_rate: Number(
          form.interest_rate,
        ),

        loan_interest: Number(
          form.loan_interest,
        ),

        pn_value: Number(
          form.pn_value,
        ),

        loan_term: Number(
          form.loan_term,
        ),

        amortization_amount:
          Number(
            form.amortization_amount,
          ),

        disbursement_date:
          form.disbursement_date,

        maturity_date:
          form.maturity_date,

        frequency:
          form.frequency ||
          "MONTHLY",

        status:
          form.status,
      };

      if (loan) {
        await updateLoan(
          loan.id,
          payload as UpdateLoanRequest,
        );
      } else {
        await createLoan(
          payload as CreateLoanRequest,
        );
      }

      onSuccess();
    } catch (err: any) {
      console.error(
        "Failed to save loan:",
        err,
      );

      setError(
        err?.response?.data?.message ||
          err?.message ||
          "Unable to save loan.",
      );
    } finally {
      setSaving(false);
    }
  };

  return (
    <Dialog
      open={open}
      onClose={
        saving ? undefined : onClose
      }
      fullWidth
      maxWidth="sm"
    >
      <DialogTitle
        sx={{
          fontWeight: 700,
        }}
      >
        {editing
          ? "Edit Loan"
          : "Add Loan"}
      </DialogTitle>

      <DialogContent>
        <Stack
          spacing={2}
          sx={{ pt: 1 }}
        >
          <TextField
            label="Client ID"
            type="number"
            fullWidth
            value={form.client_id}
            onChange={(e) =>
              change(
                "client_id",
                e.target.value,
              )
            }
          />

          <TextField
            label="PN Number"
            fullWidth
            value={form.pn_number}
            onChange={(e) =>
              change(
                "pn_number",
                e.target.value,
              )
            }
          />

          <TextField
            label="Loan Types"
            fullWidth
            value={form.loan_type}
            onChange={(e) =>
              change(
                "loan_type",
                e.target.value,
              )
            }
          />

          <TextField
            label="Principal Amount"
            type="number"
            fullWidth
            value={
              form.principal_amount
            }
            onChange={(e) =>
              change(
                "principal_amount",
                e.target.value,
              )
            }
          />

          <TextField
            label="Interest Rate (%)"
            type="number"
            fullWidth
            value={form.interest_rate}
            onChange={(e) =>
              change(
                "interest_rate",
                e.target.value,
              )
            }
          />

          <TextField
            label="Loan Interest"
            type="number"
            fullWidth
            value={form.loan_interest}
            onChange={(e) =>
              change(
                "loan_interest",
                e.target.value,
              )
            }
          />

          <TextField
            label="PN Value"
            type="number"
            fullWidth
            value={form.pn_value}
            onChange={(e) =>
              change(
                "pn_value",
                e.target.value,
              )
            }
          />

          <TextField
            label="Loan Term"
            type="number"
            fullWidth
            value={form.loan_term}
            onChange={(e) =>
              change(
                "loan_term",
                e.target.value,
              )
            }
          />

          <TextField
            label="Amortization Amount"
            type="number"
            fullWidth
            value={
              form.amortization_amount
            }
            onChange={(e) =>
              change(
                "amortization_amount",
                e.target.value,
              )
            }
          />

          <TextField
            label="Disbursement Date"
            type="date"
            fullWidth
            InputLabelProps={{
              shrink: true,
            }}
            value={
              form.disbursement_date
            }
            onChange={(e) =>
              change(
                "disbursement_date",
                e.target.value,
              )
            }
          />

          <TextField
            label="Maturity Date"
            type="date"
            fullWidth
            InputLabelProps={{
              shrink: true,
            }}
            value={
              form.maturity_date
            }
            onChange={(e) =>
              change(
                "maturity_date",
                e.target.value,
              )
            }
          />

          <TextField
            select
            label="Frequency"
            fullWidth
            value={form.frequency}
            onChange={(e) =>
              change(
                "frequency",
                e.target.value,
              )
            }
          >
            <MenuItem value="MONTHLY">
              Monthly
            </MenuItem>

            <MenuItem value="SEMI-MONTHLY">
              Semi-Monthly
            </MenuItem>
          </TextField>

          <TextField
            select
            label="Status"
            fullWidth
            value={form.status}
            onChange={(e) =>
              change(
                "status",
                e.target.value,
              )
            }
          >
            <MenuItem value="ACTIVE">
              Active
            </MenuItem>

            <MenuItem value="PAID">
              Paidzz
            </MenuItem>

            <MenuItem value="CLOSED">
              Closed
            </MenuItem>

            <MenuItem value="DEFAULTED">
              Defaulted
            </MenuItem>
          </TextField>

          {error && (
            <div
              style={{
                color: "#A51D1D",
                fontSize: 14,
              }}
            >
              {error}
            </div>
          )}
        </Stack>
      </DialogContent>

      <DialogActions
        sx={{
          px: 3,
          pb: 2.5,
        }}
      >
        <Button
          onClick={onClose}
          disabled={saving}
        >
          Cancel
        </Button>

        <Button
          variant="contained"
          onClick={handleSubmit}
          disabled={saving}
          sx={{
            bgcolor: "#8F2115",
            "&:hover": {
              bgcolor: "#70150F",
            },
          }}
        >
          {saving
            ? "Saving..."
            : editing
              ? "Update Loan"
              : "Create Loan"}
        </Button>
      </DialogActions>
    </Dialog>
  );
}