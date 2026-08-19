import {
useEffect,
useState,
} from "react";

import {
Alert,
Box,
Button,
Card,
CardContent,
Chip,
Divider,
Grid,
Stack,
Typography,
} from "@mui/material";

import {
AccountBalanceOutlined,
ArrowBack,
PersonOutlineOutlined,
RefreshOutlined,
} from "@mui/icons-material";

import {
useNavigate,
useParams,
} from "react-router-dom";

import {
getLoanProfile,
rebuildAmortization,
} from "../../services/loanService";

import MakePaymentDialog from "./MakePaymentDialog";

import type {
LoanProfileResponse,
} from "../../types/loan";

/* =========================================================
HELPERS
========================================================= */

const money = (value: number) =>
new Intl.NumberFormat(
"en-PH",
{
style: "currency",
currency: "PHP",
},
).format(value);

const formatDate = (
value?: string | null,
) => {
if (!value) {
return "—";
}

const d = new Date(value);

if (Number.isNaN(d.getTime())) {
return value;
}

return d.toLocaleDateString(
"en-PH",
{
year: "numeric",
month: "short",
day: "numeric",
},
);
};



const formatLoanStatus = (
  status?: string,
) => {
  switch (status) {
    case "ACTIVE":
      return "Active";

    case "PAID":
      return "Paid";

    case "CLOSED":
      return "Closed";

    case "DEFAULTED":
      return "Defaulted";

    default:
      return status || "Unknown";
  }
};

/* =========================================================
PAGE
========================================================= */

export default function LoanProfilePage() {

const navigate = useNavigate();

const { id } =
useParams<{
id: string;
}>();

const [
profile,
setProfile,
 ] =
useState<
LoanProfileResponse["data"] | null
>(null);

const [
loading,
setLoading,
 ] =
useState(true);

const [
error,
setError,
 ] =
useState("");

const [
paymentOpen,
setPaymentOpen,
 ] =
useState(false);

const [
rebuilding,
setRebuilding,
 ] =
useState(false);

/* =======================================================
LOAD PROFILE
======================================================= */

const load = async () => {

const loanId =
Number(id);

if (
!id ||
!Number.isInteger(
loanId,
) ||
loanId <= 0
) {

setError(
"Invalid loan ID.",
);

setLoading(false);

return;
}

try {

setLoading(true);
setError("");

const data =
await getLoanProfile(
loanId,
);

setProfile(data);

} catch (err) {

console.error(
"Failed to load loan profile:",
err,
);

setError(
"Unable to load loan profile.",
);

} finally {

setLoading(false);
}
};

useEffect(() => {
void load();
}, [id]);

/* =======================================================
REBUILD
======================================================= */

const handleRebuild =
async () => {

const loanId =
Number(id);

if (
!Number.isInteger(
loanId,
) ||
loanId <= 0
) {
return;
}

const confirmed =
window.confirm(
"Rebuild the amortization schedule?\n\nExisting payments will be preserved and reapplied to the new schedule.",
);

if (!confirmed) {
return;
}

try {

setRebuilding(true);
setError("");

await rebuildAmortization(
loanId,
);

await load();

} catch (err) {

console.error(
"Failed to rebuild amortization:",
err,
);

setError(
"Unable to rebuild amortization schedule.",
);

} finally {

setRebuilding(false);
}
};

/* =======================================================
LOADING
======================================================= */

if (loading) {

return (
<Box
sx={{
minHeight: 400,
display: "flex",
alignItems: "center",
justifyContent: "center",
}}
>
<Typography
color="#756B68"
>
Loading loan profile...
</Typography>
</Box>
);
}

/* =======================================================
ERROR
======================================================= */

if (!profile || error) {

return (
<Box>

<Button
startIcon={
<ArrowBack />
}
onClick={() =>
navigate(
"/loans",
)
}
sx={{
color: "#8F2115",
}}
>
Back to Loans
</Button>

<Alert
severity="error"
sx={{
mt: 2,
}}
>
{error ||
"Loan profile not found."}
</Alert>

</Box>
);
}

const {
loan,
client,
summary,
amortizations,
payments,
} = profile;

/* =======================================================
UI
======================================================= */

return (
<Box>

{/* ===================================================
HEADER
=================================================== */}

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
sx={{
mb: 3,
}}
>

<Box>

<Button
startIcon={
<ArrowBack />
}
onClick={() =>
navigate(
"/loans",
)
}
sx={{
px: 0,
color: "#8F2115",
}}
>
Back to Loans
</Button>

<Typography
variant="h4"
fontWeight={700}
sx={{
color: "#2B211F",
}}
>
{loan.pn_number}
</Typography>

<Typography
color="#756B68"
>
Loan Profile
</Typography>

</Box>

<Stack
direction={{
xs: "column",
sm: "row",
}}
spacing={1}
>

<Button
variant="outlined"
startIcon={
<RefreshOutlined />
}
disabled={rebuilding}
onClick={
handleRebuild
}
sx={{
borderColor:
"#8F2115",

color:
"#8F2115",

"&:hover": {
borderColor:
"#70150F",

backgroundColor:
"rgba(143,33,21,0.05)",
},
}}
>
{rebuilding
? "Rebuilding..."
: "Rebuild Amortization"}
</Button>

<Chip
label={
formatLoanStatus(loan.status)
}
sx={{
fontWeight: 700,

backgroundColor:
loan.status ===
"ACTIVE"
? "#F2E5C1"
: "#ECE8E6",

color:
loan.status ===
"ACTIVE"
? "#705B18"
: "#665D59",
}}
/>

<Button
variant="contained"
disabled={
summary.outstanding_balance <=
0 ||
rebuilding
}
onClick={() =>
setPaymentOpen(
true,
)
}
sx={{
backgroundColor:
"#8F2115",

"&:hover": {
backgroundColor:
"#70150F",
},
}}
>
Make Payment
</Button>

</Stack>

</Stack>

{/* ===================================================
ERROR
=================================================== */}

{error && (
<Alert
severity="error"
sx={{
mb: 3,
}}
onClose={() =>
setError("")
}
>
{error}
</Alert>
)}

{/* ===================================================
LOAN INFORMATION
=================================================== */}

<Card
sx={{
mb: 3,
}}
>

<CardContent>

<SectionTitle
icon={
<AccountBalanceOutlined />
}
text="Loan Information"
/>

<Divider
sx={{
mb: 3,
}}
/>

<Grid
container
spacing={3}
>

<Info
label="PN Number"
value={
loan.pn_number
}
/>

<Info
label="Loan Type"
value={
loan.loan_type ||
"—"
}
/>

<Info
label="Interest Rate"
value={`${loan.interest_rate}%`}
/>

<Info
label="Loan Term"
value={`${loan.loan_term} ${loan.frequency === "SEMI-MONTHLY" ? "months" : "months"}`}
/>

<Info
label="Frequency"
value={
loan.frequency ||
"MONTHLY"
}
/>

<Info
label="Disbursement Date"
value={formatDate(
loan.disbursement_date,
)}
/>

<Info
label="Maturity Date"
value={formatDate(
loan.maturity_date,
)}
/>

<Info
label="Principal"
value={money(
loan.principal_amount,
)}
/>

<Info
label="Loan Interest"
value={money(
loan.loan_interest,
)}
/>

<Info
label="PN Value"
value={money(
loan.pn_value,
)}
/>

<Info
label="Amortization"
value={money(
loan.amortization_amount,
)}
/>

</Grid>

</CardContent>

</Card>

{/* ===================================================
CLIENT INFORMATION
=================================================== */}

<Card
sx={{
mb: 3,
}}
>

<CardContent>

<SectionTitle
icon={
<PersonOutlineOutlined />
}
text="Client Information"
/>

<Divider
sx={{
mb: 3,
}}
/>

<Grid
container
spacing={3}
>

<Info
label="Client ID"
value={`#${client.id}`}
/>

<Info
label="Name"
value={`${client.first_name} ${client.last_name}`}
/>

<Info
label="Contact"
value={
client.contact_number ||
"—"
}
/>

<Info
label="Email"
value={
client.email ||
"—"
}
/>

<Grid
size={{
xs: 12,
}}
>

<Info
label="Address"
value={
client.current_address ||
"—"
}
/>

</Grid>

</Grid>

</CardContent>

</Card>

{/* ===================================================
SUMMARY
=================================================== */}

<Grid
container
spacing={2}
sx={{
mb: 3,
}}
>

<Summary
title="Principal"
value={
summary.principal_amount
}
/>

<Summary
title="PN Value"
value={
summary.pn_value
}
/>

<Summary
title="Total Paid"
value={
summary.total_paid
}
/>

<Summary
title="Outstanding"
value={
summary.outstanding_balance
}
/>

</Grid>

{/* ===================================================
AMORTIZATION
=================================================== */}

<Card
sx={{
mb: 3,
}}
>

<CardContent>

<Typography
variant="h6"
fontWeight={700}
sx={{
color: "#2B211F",
}}
>
Amortization Schedule
</Typography>

<Typography
variant="body2"
color="#756B68"
sx={{
mb: 2,
}}
>
Scheduled payments and
payment progress.
</Typography>

{!amortizations.length ? (

<Alert severity="info">
No amortization records.
</Alert>

) : (

<Box
sx={{
overflowX:
"auto",
}}
>

<Box
sx={{
minWidth: 1050,
}}
>

<ScheduleHeader />

{amortizations.map(
(item) => (
<ScheduleRow
key={
item.id
}
item={
item
}
/>
),
)}

</Box>

</Box>

)}

</CardContent>

</Card>

{/* ===================================================
PAYMENTS
=================================================== */}

<Card
sx={{
mb: 3,
}}
>

<CardContent>

<Typography
variant="h6"
fontWeight={700}
sx={{
color: "#2B211F",
}}
>
Payment History
</Typography>

<Typography
variant="body2"
color="#756B68"
sx={{
mb: 2,
}}
>
Payments are applied to
interest first, then
principal.
</Typography>

{!payments.length ? (

<Alert severity="info">
No payment records.
</Alert>

) : (

<Box
sx={{
overflowX:
"auto",
}}
>

<Box
sx={{
minWidth: 950,
}}
>

<PaymentHeader />

{payments.map(
(payment) => (
<PaymentRow
key={
payment.id
}
payment={
payment
}
/>
),
)}

</Box>

</Box>

)}

</CardContent>

</Card>

{/* ===================================================
PAYMENT DIALOG
=================================================== */}

<MakePaymentDialog
open={
paymentOpen
}
loanId={
loan.id
}
outstandingBalance={
summary.outstanding_balance
}
onClose={() =>
setPaymentOpen(
false,
)
}
onSuccess={async () => {

setPaymentOpen(
false,
);

await load();

}}
/>

</Box>
);
}

/* =========================================================
SECTION TITLE
========================================================= */

function SectionTitle({
icon,
text,
}: {
icon: React.ReactNode;
text: string;
}) {

return (
<Stack
direction="row"
spacing={1}
alignItems="center"
sx={{
mb: 2,
}}
>

<Box
sx={{
color: "#8F2115",
display: "flex",
}}
>
{icon}
</Box>

<Typography
variant="h6"
fontWeight={700}
sx={{
color: "#2B211F",
}}
>
{text}
</Typography>

</Stack>
);
}

/* =========================================================
INFO
========================================================= */

function Info({
label,
value,
}: {
label: string;
value: string;
}) {

return (
<Grid
size={{
xs: 12,
sm: 6,
md: 3,
}}
>

<Typography
variant="caption"
color="text.secondary"
>
{label}
</Typography>

<Typography
fontWeight={600}
sx={{
mt: 0.5,
wordBreak:
"break-word",
}}
>
{value}
</Typography>

</Grid>
);
}

/* =========================================================
SUMMARY
========================================================= */

function Summary({
title,
value,
}: {
title: string;
value: number;
}) {

return (
<Grid
size={{
xs: 12,
sm: 6,
md: 3,
}}
>

<Card
sx={{
height: "100%",
}}
>

<CardContent>

<Typography
color="text.secondary"
variant="body2"
>
{title}
</Typography>

<Typography
variant="h6"
fontWeight={700}
sx={{
mt: 1,
color: "#2B211F",
}}
>
{money(value)}
</Typography>

</CardContent>

</Card>

</Grid>
);
}

/* =========================================================
AMORTIZATION HEADER
========================================================= */

function ScheduleHeader() {

// const headers = [
// "Due Date",
// "Principal",
// "Interest",
// "Total",
// "Paid Principal",
// "Paid Interest",
// "Status",
//  ];

return (
<Grid
container
sx={{
px: 2,
py: 1.5,
bgcolor: "#F8F5F3",
borderBottom:
"1px solid #E7DDD9",
fontWeight: 700,
alignItems:
"center",
}}
>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
Due Date
</Grid>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
Principal
</Grid>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
Interest
</Grid>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
Total
</Grid>

<Grid
size={{ xs: 12, sm: 1.8 }}
>
Paid Principal
</Grid>

<Grid
size={{ xs: 12, sm: 1.8 }}
>
Paid Interest
</Grid>

<Grid
size={{ xs: 12, sm: 2.4 }}
>
Status
</Grid>

</Grid>
);
}

/* =========================================================
AMORTIZATION ROW
========================================================= */

function ScheduleRow({
item,
}: {
item: {
due_date: string;
principal_amount: number;
interest_amount: number;
total_amount: number;
paid_principal_amount?: number;
paid_interest_amount?: number;
status: string;
};
}) {

return (
<Grid
container
sx={{
px: 2,
py: 1.5,
borderBottom:
"1px solid #F0EAE7",
alignItems:
"center",
}}
>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
{formatDate(
item.due_date,
)}
</Grid>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
{money(
item.principal_amount,
)}
</Grid>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
{money(
item.interest_amount,
)}
</Grid>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
<b>
{money(
item.total_amount,
)}
</b>
</Grid>

<Grid
size={{ xs: 12, sm: 1.8 }}
>
{money(
item.paid_principal_amount ??
0,
)}
</Grid>

<Grid
size={{ xs: 12, sm: 1.8 }}
>
{money(
item.paid_interest_amount ??
0,
)}
</Grid>

<Grid
size={{ xs: 12, sm: 2.4 }}
>

<Chip
size="small"
label={
item.status
}
sx={{
fontWeight: 700,

backgroundColor:
item.status ===
"PAID"
? "#E4F0E5"
: item.status ===
"PARTIAL"
? "#F2E5C1"
: "#ECE8E6",

color:
item.status ===
"PAID"
? "#35633B"
: item.status ===
"PARTIAL"
? "#705B18"
: "#665D59",
}}
/>

</Grid>

</Grid>
);
}

/* =========================================================
PAYMENT HEADER
========================================================= */

function PaymentHeader() {

return (
<Grid
container
sx={{
px: 2,
py: 1.5,
bgcolor: "#F8F5F3",
borderBottom:
"1px solid #E7DDD9",
fontWeight: 700,
}}
>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
Date
</Grid>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
Amount
</Grid>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
Channel
</Grid>

<Grid
size={{ xs: 12, sm: 2 }}
>
Reference
</Grid>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
Interest
</Grid>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
Principal
</Grid>

<Grid
size={{ xs: 12, sm: 2.5 }}
>
Outstanding
</Grid>

</Grid>
);
}

/* =========================================================
PAYMENT ROW
========================================================= */

function PaymentRow({
payment,
}: {
payment: {
payment_date: string;
amount_paid: number;
payment_channel?: string | null;
reference_number?: string | null;
interest_applied: number;
principal_applied: number;
outstanding_balance: number;
};
}) {

return (
<Grid
container
sx={{
px: 2,
py: 1.5,
borderBottom:
"1px solid #F0EAE7",
alignItems:
"center",
}}
>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
{formatDate(
payment.payment_date,
)}
</Grid>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
<b>
{money(
payment.amount_paid,
)}
</b>
</Grid>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
{payment.payment_channel ||
"—"}
</Grid>

<Grid
size={{ xs: 12, sm: 2 }}
>
{payment.reference_number ||
"—"}
</Grid>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
{money(
payment.interest_applied,
)}
</Grid>

<Grid
size={{ xs: 12, sm: 1.5 }}
>
{money(
payment.principal_applied,
)}
</Grid>

<Grid
size={{ xs: 12, sm: 2.5 }}
>
<b>
{money(
payment.outstanding_balance,
)}
</b>
</Grid>

</Grid>
);
}