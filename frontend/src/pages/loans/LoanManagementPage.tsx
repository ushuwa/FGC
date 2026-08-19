import { useEffect, useState } from "react";

import {
Alert,
Box,
Button,
Card,
CardContent,
Dialog,
DialogActions,
DialogContent,
DialogTitle,
FormControl,
Grid,
IconButton,
InputAdornment,
InputLabel,
MenuItem,
Select,
Stack,
TextField,
Typography,
} from "@mui/material";

import {
Add,
EditOutlined,
Search,
VisibilityOutlined,
} from "@mui/icons-material";

import {
DataGrid,
type GridColDef,
} from "@mui/x-data-grid";

import { useNavigate } from "react-router-dom";

import {
createLoan,
getLoans,
updateLoan,
} from "../../services/loanService";

import { getClients } from "../../services/clientService";

import type {
Loan,
CreateLoanRequest,
} from "../../types/loan";

import type { Client } from "../../types/client";

interface LoanFormValues {
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

const emptyForm: LoanFormValues = {
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

export default function LoanManagementPage() {
const navigate = useNavigate();

const [loans, setLoans] = useState<Loan[]>([]);
const [clients, setClients] = useState<Client[]>([]);

const [search, setSearch] = useState("");
const [status, setStatus] = useState("");

const [loading, setLoading] = useState(true);
const [saving, setSaving] = useState(false);

const [error, setError] = useState("");

const [formOpen, setFormOpen] = useState(false);
const [editingLoan, setEditingLoan] =
useState<Loan | null>(null);

const [form, setForm] =
useState<LoanFormValues>(emptyForm);

// ============================================
// LOAD LOANS
// ============================================

const loadLoans = async (
searchValue = search,
statusValue = status,
) => {
try {
setLoading(true);
setError("");

const data = await getLoans(
searchValue,
statusValue,
);

setLoans(data);
} catch (err) {
console.error(
"Failed to load loans:",
err,
);

setError(
"Unable to load loans.",
);
} finally {
setLoading(false);
}
};

// ============================================
// LOAD CLIENTS
// ============================================

const loadClients = async () => {
try {
const data = await getClients("");
setClients(data);
} catch (err) {
console.error(
"Failed to load clients:",
err,
);
}
};

useEffect(() => {
void loadLoans("", "");
void loadClients();
}, []);

// ============================================
// SEARCH
// ============================================

const handleSearch = () => {
void loadLoans(
search,
status,
);
};

// ============================================
// FORM
// ============================================

const updateField = (
  field: keyof LoanFormValues,
  value: string,
) => {
  setForm((current) => {
    const updated = {
      ...current,
      [field]: value,
    };

    // Automatically calculate PN Value
    if (
      field === "principal_amount" ||
      field === "loan_interest"
    ) {
      updated.pn_value = calculatePNValue(
        updated.principal_amount,
        updated.loan_interest,
      );
    }

    // Automatically calculate Amortization
    if (
      field === "principal_amount" ||
      field === "loan_interest" ||
      field === "loan_term" ||
      field === "frequency"
    ) {
      updated.amortization_amount =
        calculateAmortization(
          updated.principal_amount,
          updated.loan_interest,
          updated.loan_term,
          updated.frequency,
        );
    }

    return updated;
  });
};



const calculateAmortization = (
  principal: string,
  loanInterest: string,
  loanTerm: string,
  frequency: string,
) => {
  const principalAmount = Number(principal);
  const interestAmount = Number(loanInterest);
  const term = Number(loanTerm);

  if (
    principalAmount <= 0 ||
    interestAmount < 0 ||
    term <= 0
  ) {
    return "";
  }

  const periods =
    frequency === "SEMI-MONTHLY"
      ? term * 2
      : term;

  if (periods <= 0) {
    return "";
  }

  return (
    (principalAmount + interestAmount) /
    periods
  ).toFixed(2);
};


useEffect(() => {
  const amortization =
    calculateAmortization(
      form.principal_amount,
      form.loan_interest,
      form.loan_term,
      form.frequency,
    );

  setForm((current) => {
    if (
      current.amortization_amount ===
      amortization
    ) {
      return current;
    }

    return {
      ...current,
      amortization_amount:
        amortization,
    };
  });
}, [
  form.principal_amount,
  form.loan_interest,
  form.loan_term,
  form.frequency,
]);

const calculatePNValue = (
  principal: string,
  loanInterest: string,
) => {
  const principalAmount = Number(principal);
  const interestAmount = Number(loanInterest);

  if (
    principalAmount <= 0 ||
    interestAmount < 0
  ) {
    return "";
  }

  return (
    principalAmount + interestAmount
  ).toFixed(2);
};




const calculateMaturityDate = (
  disbursementDate: string,
  loanTerm: string,
  frequency: string,
) => {
  if (
    !disbursementDate ||
    Number(loanTerm) <= 0
  ) {
    return "";
  }

  const date = new Date(
    `${disbursementDate}T00:00:00`,
  );

  const term = Number(loanTerm);

  if (frequency === "MONTHLY") {
    date.setMonth(
      date.getMonth() + term,
    );

    return date
      .toISOString()
      .slice(0, 10);
  }

  if (
    frequency === "SEMI-MONTHLY"
  ) {
    let currentDate = new Date(date);

    for (
      let i = 0;
      i < term * 2;
      i++
    ) {
      const day =
        currentDate.getDate();

      if (day < 15) {
        currentDate = new Date(
          currentDate.getFullYear(),
          currentDate.getMonth(),
          15,
        );
      } else if (day === 15) {
        currentDate = new Date(
          currentDate.getFullYear(),
          currentDate.getMonth() + 1,
          0,
        );
      } else {
        currentDate = new Date(
          currentDate.getFullYear(),
          currentDate.getMonth() + 1,
          15,
        );
      }
    }

    return currentDate
      .toISOString()
      .slice(0, 10);
  }

  return "";
};


useEffect(() => {
  const maturityDate =
    calculateMaturityDate(
      form.disbursement_date,
      form.loan_term,
      form.frequency,
    );

  setForm((current) => {
    if (
      current.maturity_date ===
      maturityDate
    ) {
      return current;
    }

    return {
      ...current,
      maturity_date: maturityDate,
    };
  });
}, [
  form.disbursement_date,
  form.loan_term,
  form.frequency,
]);


const handleAddLoan = () => {
setEditingLoan(null);
setForm(emptyForm);
setFormOpen(true);
};

const handleEditLoan = (loan: Loan) => {
  const frequency =
    loan.frequency || "MONTHLY";

  const principal =
    String(loan.principal_amount ?? "");

  const loanInterest =
    String(loan.loan_interest ?? "");

  const loanTerm =
    String(loan.loan_term ?? "");

  const pnValue =
    calculatePNValue(
      principal,
      loanInterest,
    );

  const amortizationAmount =
    calculateAmortization(
      principal,
      loanInterest,
      loanTerm,
      frequency,
    );

  setForm({
    client_id: String(
      loan.client_id ?? "",
    ),

    pn_number:
      loan.pn_number ?? "",

    loan_type:
      loan.loan_type ?? "",

    principal_amount:
      principal,

    interest_rate:
      String(
        loan.interest_rate ?? "",
      ),

    loan_interest:
      loanInterest,

    pn_value:
      pnValue,

    loan_term:
      loanTerm,

    amortization_amount:
      amortizationAmount,

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

    frequency,

    status:
      loan.status || "ACTIVE",
  });

  setEditingLoan(loan);
  setFormOpen(true);
};

const handleCloseForm = () => {
if (saving) {
return;
}

setFormOpen(false);
setEditingLoan(null);
setForm(emptyForm);
};

// ============================================
// SAVE
// ============================================

const handleSave = async () => {
try {
setSaving(true);
setError("");

if (!form.client_id) {
setError(
"Client is required.",
);
return;
}

if (!form.pn_number.trim()) {
setError(
"PN number is required.",
);
return;
}

if (
Number(form.principal_amount) <= 0
) {
setError(
"Principal amount must be greater than zero.",
);
return;
}

if (
Number(form.loan_term) <= 0
) {
setError(
"Loan term must be greater than zero.",
);
return;
}

const payload: CreateLoanRequest = {
client_id:
Number(form.client_id),

pn_number:
form.pn_number.trim(),

loan_type:
form.loan_type.trim() ||
null,

principal_amount:
Number(
form.principal_amount,
),

interest_rate:
Number(
form.interest_rate,
),

loan_interest:
Number(
form.loan_interest,
),

pn_value:
Number(form.pn_value),

loan_term:
Number(form.loan_term),

amortization_amount:
Number(
form.amortization_amount,
),

disbursement_date:
form.disbursement_date,

maturity_date:
form.maturity_date ||
null,

frequency:
form.frequency ||
"MONTHLY",

status:
form.status,
};

if (editingLoan) {
await updateLoan(
editingLoan.id,
payload,
);
} else {
await createLoan(payload);
}

handleCloseForm();

await loadLoans(
search,
status,
);
} catch (err: any) {
console.error(
"Failed to save loan:",
err,
);

setError(
err?.response?.data?.message ||
(
editingLoan
? "Unable to update loan."
: "Unable to create loan."
),
);
} finally {
setSaving(false);
}
};

// ============================================
// FORMAT
// ============================================

const formatCurrency = (
value: number,
) => {
return new Intl.NumberFormat(
"en-PH",
{
style: "currency",
currency: "PHP",
maximumFractionDigits: 2,
},
).format(value);
};

// ============================================
// TABLE
// ============================================

const columns: GridColDef[] = [
{
field: "pn_number",
headerName: "PN Number",
minWidth: 150,
flex: 0.8,
},

{
field: "client_name",
headerName: "Client",
minWidth: 220,
flex: 1.2,
},

{
field: "loan_type",
headerName: "Loan Type",
minWidth: 150,
flex: 0.9,

renderCell: (params) =>
params.value || "—",
},

{
field: "principal_amount",
headerName: "Principal",
minWidth: 150,
flex: 0.8,

renderCell: (params) =>
formatCurrency(
Number(params.value),
),
},

{
field: "pn_value",
headerName: "PN Value",
minWidth: 150,
flex: 0.8,

renderCell: (params) =>
formatCurrency(
Number(params.value),
),
},

{
field: "total_paid",
headerName: "Paid",
minWidth: 140,
flex: 0.8,

renderCell: (params) =>
formatCurrency(
Number(params.value),
),
},

{
field: "outstanding_balance",
headerName: "Outstanding",
minWidth: 160,
flex: 0.9,

renderCell: (params) =>
formatCurrency(
Number(params.value),
),
},

{
field: "status",
headerName: "Status",
minWidth: 120,
flex: 0.6,
},

{
field: "actions",
headerName: "Actions",
width: 130,
sortable: false,
filterable: false,

renderCell: (params) => (
<Stack
direction="row"
spacing={0.5}
>
<IconButton
size="small"
onClick={() =>
navigate(`/loans/${params.row.id}`)
}
sx={{
color: "#8F2115",
}}
>
<VisibilityOutlined
fontSize="small"
/>
</IconButton>

<IconButton
size="small"
onClick={() =>
handleEditLoan(
params.row,
)
}
sx={{
color: "#D0B050",
}}
>
<EditOutlined
fontSize="small"
/>
</IconButton>
</Stack>
),
},
];

// ============================================
// UI
// ============================================

return (
<Box sx={{ width: "100%" }}>
{/* HEADER */}

<Stack
direction={{
xs: "column",
sm: "row",
}}
justifyContent="space-between"
alignItems={{
xs: "flex-start",
sm: "center",
}}
spacing={2}
sx={{ mb: 3 }}
>
<Box>
<Typography
variant="h4"
fontWeight={700}
sx={{
color: "#2B211F",
}}
>
Loan Management
</Typography>

<Typography
variant="body2"
sx={{
mt: 0.5,
color: "#756B68",
}}
>
Manage and monitor
client loans.
</Typography>
</Box>

<Button
variant="contained"
startIcon={<Add />}
onClick={handleAddLoan}
sx={{
backgroundColor: "#8F2115",

"&:hover": {
backgroundColor: "#70150F",
},
}}
>
Add Loan
</Button>
</Stack>

{/* ERROR */}

{error && (
<Alert
severity="error"
sx={{ mb: 2 }}
onClose={() =>
setError("")
}
>
{error}
</Alert>
)}

{/* TABLE */}

<Card>
<CardContent>
{/* FILTERS */}

<Stack
direction={{
xs: "column",
md: "row",
}}
spacing={1.5}
sx={{ mb: 2.5 }}
>
<TextField
fullWidth
value={search}
placeholder="Search PN number, client, or loan type"
onChange={(event) =>
setSearch(
event.target.value,
)
}
onKeyDown={(event) => {
if (
event.key ===
"Enter"
) {
handleSearch();
}
}}
InputProps={{
startAdornment: (
<InputAdornment position="start">
<Search
sx={{
color:
"#8F2115",
}}
/>
</InputAdornment>
),
}}
/>

<FormControl
sx={{
minWidth: 180,
}}
>
<InputLabel>
Status
</InputLabel>

<Select
value={status}
label="Status"
onChange={(event) =>
setStatus(
event.target.value,
)
}
>
<MenuItem value="">
All Statuses
</MenuItem>

<MenuItem value="ACTIVE">
Active
</MenuItem>

<MenuItem value="PAID">
Paid
</MenuItem>

<MenuItem value="CLOSED">
Closed
</MenuItem>

<MenuItem value="DEFAULTED">
Defaulted
</MenuItem>
</Select>
</FormControl>

<Button
variant="outlined"
onClick={handleSearch}
sx={{
minWidth: 110,
borderColor:
"#8F2115",
color: "#8F2115",

"&:hover": {
borderColor:
"#70150F",
backgroundColor:
"rgba(143,33,21,0.05)",
},
}}
>
Search
</Button>
</Stack>

{/* DATA GRID */}

<DataGrid
rows={loans}
columns={columns}
loading={loading}
autoHeight
disableRowSelectionOnClick
pageSizeOptions={[
10,
25,
50,
 ]}
initialState={{
pagination: {
paginationModel: {
pageSize: 10,
page: 0,
},
},
}}
sx={{
border: "none",

"& .MuiDataGrid-columnHeaders":
{
backgroundColor:
"#F8F5F3",

borderBottom:
"1px solid #E7DDD9",
},

"& .MuiDataGrid-columnHeaderTitle":
{
fontWeight: 700,
color:
"#2B211F",
},

"& .MuiDataGrid-row:hover":
{
backgroundColor:
"rgba(208,176,80,0.06)",
},

"& .MuiDataGrid-cell":
{
borderBottom:
"1px solid #F0EAE7",
},
}}
/>
</CardContent>
</Card>

{/* ========================================
ADD / EDIT LOAN DIALOG
======================================== */}

<Dialog
open={formOpen}
onClose={handleCloseForm}
fullWidth
maxWidth="md"
>
<DialogTitle
sx={{
fontWeight: 700,
color: "#2B211F",
}}
>
{editingLoan
? "Edit Loan"
: "Add Loan"}
</DialogTitle>

<DialogContent>
<Stack
spacing={2}
sx={{ pt: 1 }}
>
{/* CLIENT */}

<FormControl fullWidth>
<InputLabel>
Client
</InputLabel>

<Select
value={
form.client_id
}
label="Client"
onChange={(event) =>
updateField(
"client_id",
event.target.value,
)
}
>
<MenuItem value="">
Select Client
</MenuItem>

{clients.map(
(client) => (
<MenuItem
key={client.id}
value={client.id}
>
{
client.first_name
}{" "}
{
client.last_name
}
</MenuItem>
),
)}
</Select>
</FormControl>

{/* PN + TYPE */}

<Grid
container
spacing={2}
>
<Grid
size={{
xs: 12,
sm: 6,
}}
>
<TextField
fullWidth
label="PN Number"
value={
form.pn_number
}
onChange={(event) =>
updateField(
"pn_number",
event.target.value,
)
}
/>
</Grid>

<Grid
size={{
xs: 12,
sm: 6,
}}
>
<TextField
fullWidth
label="Loan Type"
value={
form.loan_type
}
onChange={(event) =>
updateField(
"loan_type",
event.target.value,
)
}
/>
</Grid>
</Grid>

{/* FINANCIAL */}

<Grid
container
spacing={2}
>
<Grid
size={{
xs: 12,
sm: 4,
}}
>
<TextField
fullWidth
type="number"
label="Principal Amount"
value={
form.principal_amount
}
onChange={(event) =>
updateField(
"principal_amount",
event.target.value,
)
}
/>
</Grid>

<Grid
size={{
xs: 12,
sm: 4,
}}
>
<TextField
fullWidth
type="number"
label="Interest Rate (%)"
value={
form.interest_rate
}
onChange={(event) =>
updateField(
"interest_rate",
event.target.value,
)
}
/>
</Grid>

<Grid
size={{
xs: 12,
sm: 4,
}}
>
<TextField
fullWidth
type="number"
label="Loan Interest"
value={
form.loan_interest
}
onChange={(event) =>
updateField(
"loan_interest",
event.target.value,
)
}
/>
</Grid>
</Grid>

{/* PN VALUE / TERM / AMORTIZATION */}

<Grid
container
spacing={2}
>
<Grid
size={{
xs: 12,
sm: 4,
}}
>
<TextField
fullWidth
type="number"
label="PN Value"
value={
form.pn_value
}
InputProps={{
readOnly: true,
}}  
/>
</Grid>

<Grid
size={{
xs: 12,
sm: 4,
}}
>
<TextField
fullWidth
type="number"
label="Loan Term (Months)"
value={
form.loan_term
}
onChange={(event) =>
updateField(
"loan_term",
event.target.value,
)
}
/>
</Grid>

<Grid
size={{
xs: 12,
sm: 4,
}}
>
<TextField
fullWidth
type="number"
label="Amortization Amount"
value={
form.amortization_amount
}
InputProps={{
readOnly: true,
}}
// onChange={(event) =>
// updateField(
// "amortization_amount",
// event.target.value,
// )
// }
/>
</Grid>
</Grid>

{/* DATES */}

<Grid
container
spacing={2}
>
<Grid
size={{
xs: 12,
sm: 6,
}}
>
<TextField
fullWidth
type="date"
label="Disbursement Date"
value={
form.disbursement_date
}
onChange={(event) =>
updateField(
"disbursement_date",
event.target.value,
)
}
slotProps={{
inputLabel: {
shrink: true,
},
}}
/>
</Grid>

<Grid
size={{
xs: 12,
sm: 6,
}}
>
<TextField
fullWidth
type="date"
label="Maturity Date"
value={
form.maturity_date
}
// onChange={(event) =>
// updateField(
// "maturity_date",
// event.target.value,
// )
// }
InputProps={{
readOnly: true,
}}
slotProps={{
inputLabel: {
shrink: true,
},
}}
/>
</Grid>
</Grid>

{/* FREQUENCY */}

<FormControl fullWidth>
<InputLabel>
Frequency
</InputLabel>

<Select
value={
form.frequency
}
label="Frequency"
onChange={(event) =>
updateField(
"frequency",
event.target.value,
)
}
>
<MenuItem value="MONTHLY">
Monthly
</MenuItem>

<MenuItem value="SEMI-MONTHLY">
Semi-Monthly
</MenuItem>
</Select>
</FormControl>

{/* STATUS */}

<FormControl fullWidth>
<InputLabel>
Status
</InputLabel>

<Select
value={form.status}
label="Status"
onChange={(event) =>
updateField(
"status",
event.target.value,
)
}
>
<MenuItem value="ACTIVE">
Active
</MenuItem>

<MenuItem value="PAID">
Paid
</MenuItem>

<MenuItem value="CLOSED">
Closed
</MenuItem>

<MenuItem value="DEFAULTED">
Defaulted
</MenuItem>
</Select>
</FormControl>
</Stack>
</DialogContent>

<DialogActions
sx={{
px: 3,
pb: 2.5,
}}
>
<Button
onClick={handleCloseForm}
disabled={saving}
sx={{
color: "#756B68",
}}
>
Cancel
</Button>

<Button
variant="contained"
onClick={handleSave}
disabled={saving}
sx={{
minWidth: 120,
backgroundColor:
"#8F2115",

"&:hover": {
backgroundColor:
"#70150F",
},
}}
>
{saving
? "Saving..."
: editingLoan
? "Update Loan"
: "Save Loan"}
</Button>
</DialogActions>
</Dialog>
</Box>
);
}