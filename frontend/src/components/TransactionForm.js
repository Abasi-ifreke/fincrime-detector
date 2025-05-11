import React, { useState } from 'react';

function TransactionForm({ onSubmit }) {
  const [transaction, setTransaction] = useState({
    account_id: '',
    transaction_type: '',
    amount: '',
  });
  const [amountError, setAmountError] = useState('');

  const handleChange = (event) => {
    const { name, value } = event.target;
    setTransaction({ ...transaction, [name]: value });
    if (name === 'amount') {
      setAmountError(''); // Clear error on input change
    }
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    const parsedAmount = parseFloat(transaction.amount);
    if (isNaN(parsedAmount)) {
      setAmountError('Amount must be a valid number.');
      return;
    }

    const transactionData = {
      account_id: transaction.account_id,
      transaction_type: transaction.transaction_type,
      amount: parsedAmount, // Ensure amount is a number
    };

    console.log('Submitting transaction data:', transactionData);
    onSubmit(transactionData);
    setTransaction({ account_id: '', transaction_type: '', amount: '' }); // Clear form
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label htmlFor="account_id">Account ID:</label>
        <input
          type="text"
          id="account_id"
          name="account_id"
          value={transaction.account_id}
          onChange={handleChange}
          required
        />
      </div>
      <div>
        <label htmlFor="transaction_type">Transaction Type:</label>
        <input
          type="text"
          id="transaction_type"
          name="transaction_type"
          value={transaction.transaction_type}
          onChange={handleChange}
          required
        />
      </div>
      <div>
        <label htmlFor="amount">Amount:</label>
        <input
          type="number"
          id="amount"
          name="amount"
          value={transaction.amount}
          onChange={handleChange}
          required
        />
        {amountError && <p className="error-message">{amountError}</p>}
      </div>
      <button type="submit">Submit Transaction</button>
    </form>
  );
}

export default TransactionForm;