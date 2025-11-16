use borsh::{BorshDeserialize, BorshSerialize};
use solana_program::{
    account_info::{next_account_info, AccountInfo},
    entrypoint,
    entrypoint::ProgramResult,
    msg,
    program_error::ProgramError,
    pubkey::Pubkey,
};

// Trade data structure
#[derive(BorshSerialize, BorshDeserialize, Debug)]
pub struct TradeRecord {
    pub user_id: [u8; 16],      // UUID as bytes
    pub asset_symbol: String,    // Asset symbol (BTC, ETH, etc.)
    pub trade_type: TradeType,   // BUY or SELL
    pub quantity: u64,           // Quantity * 10^8 for precision
    pub price: u64,              // Price * 10^2 for precision
    pub timestamp: i64,          // Unix timestamp
}

#[derive(BorshSerialize, BorshDeserialize, Debug, PartialEq)]
pub enum TradeType {
    Buy,
    Sell,
}

// Instruction enum
#[derive(BorshSerialize, BorshDeserialize, Debug)]
pub enum TradeInstruction {
    /// Record a new trade
    /// Accounts expected:
    /// 0. `[signer]` The account of the user placing the trade
    /// 1. `[writable]` The trade record account to write to
    RecordTrade {
        user_id: [u8; 16],
        asset_symbol: String,
        trade_type: TradeType,
        quantity: u64,
        price: u64,
        timestamp: i64,
    },
}

// Program entrypoint
entrypoint!(process_instruction);

// Program logic
pub fn process_instruction(
    program_id: &Pubkey,
    accounts: &[AccountInfo],
    instruction_data: &[u8],
) -> ProgramResult {
    let instruction = TradeInstruction::try_from_slice(instruction_data)
        .map_err(|_| ProgramError::InvalidInstructionData)?;

    match instruction {
        TradeInstruction::RecordTrade {
            user_id,
            asset_symbol,
            trade_type,
            quantity,
            price,
            timestamp,
        } => {
            msg!("Instruction: RecordTrade");
            record_trade(
                program_id,
                accounts,
                user_id,
                asset_symbol,
                trade_type,
                quantity,
                price,
                timestamp,
            )
        }
    }
}

fn record_trade(
    _program_id: &Pubkey,
    accounts: &[AccountInfo],
    user_id: [u8; 16],
    asset_symbol: String,
    trade_type: TradeType,
    quantity: u64,
    price: u64,
    timestamp: i64,
) -> ProgramResult {
    let accounts_iter = &mut accounts.iter();
    let user_account = next_account_info(accounts_iter)?;
    let trade_account = next_account_info(accounts_iter)?;

    // Verify the user is a signer
    if !user_account.is_signer {
        return Err(ProgramError::MissingRequiredSignature);
    }

    // Create trade record
    let trade_record = TradeRecord {
        user_id,
        asset_symbol: asset_symbol.clone(),
        trade_type,
        quantity,
        price,
        timestamp,
    };

    // Serialize and store
    trade_record
        .serialize(&mut &mut trade_account.data.borrow_mut()[..])
        .map_err(|_| ProgramError::AccountDataTooSmall)?;

    msg!(
        "Trade recorded: User {:?}, Asset: {}, Type: {:?}, Quantity: {}, Price: {}, Timestamp: {}",
        user_id,
        asset_symbol,
        trade_record.trade_type,
        quantity,
        price,
        timestamp
    );

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    use solana_program::clock::Epoch;
    use solana_sdk::signature::{Keypair, Signer};

    #[test]
    fn test_record_trade() {
        let program_id = Pubkey::new_unique();
        let user_keypair = Keypair::new();
        let user_pubkey = user_keypair.pubkey();
        let trade_pubkey = Pubkey::new_unique();

        let mut user_lamports = 0;
        let mut trade_lamports = 0;
        let mut trade_data = vec![0; 256];

        let user_account = AccountInfo::new(
            &user_pubkey,
            true,
            false,
            &mut user_lamports,
            &mut [],
            &program_id,
            false,
            Epoch::default(),
        );

        let trade_account = AccountInfo::new(
            &trade_pubkey,
            false,
            true,
            &mut trade_lamports,
            &mut trade_data,
            &program_id,
            false,
            Epoch::default(),
        );

        let accounts = vec![user_account, trade_account];

        let user_id = [1u8; 16];
        let asset_symbol = "BTC".to_string();
        let trade_type = TradeType::Buy;
        let quantity = 100000000; // 1 BTC
        let price = 4500000; // $45,000
        let timestamp = 1234567890;

        let result = record_trade(
            &program_id,
            &accounts,
            user_id,
            asset_symbol,
            trade_type,
            quantity,
            price,
            timestamp,
        );

        assert!(result.is_ok());
    }
}
