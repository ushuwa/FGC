import { useState } from "react";

import {
  Alert,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Stack,
  TextField,
} from "@mui/material";

import {
  createLoanPayment,
} from "../../services/loanService";

import type {
  CreatePaymentRequest,
} from "../../types/loan";

interface Props {
  open: boolean;
  loanId: number;
  outstandingBalance: number;
  onClose: () => void;
  onSuccess: () => void | Promise<void>;
}

export default function MakePaymentDialog({
  open,
  loanId,
  outstandingBalance,
  onClose,
  onSuccess,
}: Props) {
  const [amount, setAmount] = useState("");
  const [date, setDate] = useState(
    new Date().toISOString().split("T")[0],
  );
  const [channel, setChannel] = useState("");
  const [reference, setReference] = useState("");

  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  const reset = () => {
    setAmount("");
    setDate(
      new Date().toISOString().split("T")[0],
    );
    setChannel("");
    setReference("");
    setError("");
  };

  const handleClose = () => {
    if (saving) return;

    reset();
    onClose();
  };

  const handleSubmit = async () => {
    const paymentAmount = Number(amount);

    if (!paymentAmount || paymentAmount <= 0) {
      setError("Enter a valid payment amount.");
      return;
    }

    if (paymentAmount > outstandingBalance) {
      setError(
        "Payment cannot be greater than the outstanding balance.",
      );
      return;
    }

    if (!date) {
      setError("Payment date is required.");
      return;
    }

    try {
      setSaving(true);
      setError("");

      const payload: CreatePaymentRequest = {
        payment_date: date,
        amount_paid: paymentAmount,
        payment_channel:
          channel.trim() || null,
        reference_number:
          reference.trim() || null,
      };

      

      await createLoanPayment(
        loanId,
        payload,
      );

      reset();

      await onSuccess();
    } catch (err) {
      console.error(
        "Failed to create payment:",
        err,
      );

      setError(
        "Unable to record payment.",
      );
    } finally {
      setSaving(false);
    }
  };

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      fullWidth
      maxWidth="sm"
    >
      <DialogTitle
        sx={{
          fontWeight: 700,
          color: "#2B211F",
        }}
      >
        Make Payment
      </DialogTitle>

      <DialogContent>
        <Stack
          spacing={2}
          sx={{ pt: 1 }}
        >
          {error && (
            <Alert severity="error">
              {error}
            </Alert>
          )}

          <Alert severity="info">
            Outstanding balance:{" "}
            <strong>
              ₱
              {outstandingBalance.toLocaleString(
                "en-PH",
                {
                  minimumFractionDigits: 2,
                },
              )}
            </strong>
          </Alert>

          <TextField
            label="Payment Date"
            type="date"
            value={date}
            onChange={(e) =>
              setDate(e.target.value)
            }
            fullWidth
            InputLabelProps={{
              shrink: true,
            }}
          />

          <TextField
            label="Payment Amount"
            type="number"
            value={amount}
            onChange={(e) =>
              setAmount(e.target.value)
            }
            fullWidth
            inputProps={{
              min: 0,
              max: outstandingBalance,
              step: "0.01",
            }}
          />

          <TextField
            label="Payment Channel"
            placeholder="Cash, Bank Transfer, GCash..."
            value={channel}
            onChange={(e) =>
              setChannel(e.target.value)
            }
            fullWidth
          />

          <TextField
            label="Reference Number"
            value={reference}
            onChange={(e) =>
              setReference(e.target.value)
            }
            fullWidth
          />
        </Stack>
      </DialogContent>

      <DialogActions
        sx={{
          px: 3,
          pb: 2.5,
        }}
      >
        <Button
          onClick={handleClose}
          disabled={saving}
          sx={{
            color: "#756B68",
          }}
        >
          Cancel
        </Button>

        <Button
          variant="contained"
          onClick={handleSubmit}
          disabled={saving}
          sx={{
            backgroundColor: "#8F2115",

            "&:hover": {
              backgroundColor: "#70150F",
            },
          }}
        >
          {saving
            ? "Saving..."
            : "Record Payment"}
        </Button>
      </DialogActions>
    </Dialog>
  );
}