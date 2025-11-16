import React from 'react'
import { useQuery } from 'react-query'
import { tradeService } from '../services/api'

const History = () => {
  const { data: trades, isLoading } = useQuery('tradeHistory', tradeService.getHistory)

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center text-gray-400">Loading...</div>
      </div>
    )
  }

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Trade History</h1>

      <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
        {!trades || trades.length === 0 ? (
          <div className="text-center text-gray-400 py-12">
            <p className="text-lg">No trades yet</p>
            <p className="text-sm mt-2">Start trading to see your history</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="text-left border-b border-slate-700">
                  <th className="pb-3 text-gray-400 font-medium">Date</th>
                  <th className="pb-3 text-gray-400 font-medium">Type</th>
                  <th className="pb-3 text-gray-400 font-medium">Asset</th>
                  <th className="pb-3 text-gray-400 font-medium text-right">Quantity</th>
                  <th className="pb-3 text-gray-400 font-medium text-right">Price</th>
                  <th className="pb-3 text-gray-400 font-medium text-right">Total</th>
                  <th className="pb-3 text-gray-400 font-medium">Blockchain</th>
                </tr>
              </thead>
              <tbody>
                {trades.map((trade) => (
                  <tr key={trade.id} className="border-b border-slate-700 hover:bg-slate-750">
                    <td className="py-4 text-sm">{formatDate(trade.created_at)}</td>
                    <td className="py-4">
                      <span
                        className={`px-3 py-1 rounded-full text-xs font-medium ${
                          trade.trade_type === 'BUY'
                            ? 'bg-green-500 bg-opacity-20 text-green-400'
                            : 'bg-red-500 bg-opacity-20 text-red-400'
                        }`}
                      >
                        {trade.trade_type}
                      </span>
                    </td>
                    <td className="py-4">
                      <div>
                        <div className="font-medium">{trade.asset?.symbol}</div>
                        <div className="text-sm text-gray-400">{trade.asset?.name}</div>
                      </div>
                    </td>
                    <td className="py-4 text-right">{trade.quantity?.toFixed(4)}</td>
                    <td className="py-4 text-right">${trade.price?.toFixed(2)}</td>
                    <td className="py-4 text-right font-semibold">
                      ${trade.total_amount?.toFixed(2)}
                    </td>
                    <td className="py-4">
                      {trade.solana_signature ? (
                        <a
                          href={`https://explorer.solana.com/tx/${trade.solana_signature}?cluster=devnet`}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-blue-400 hover:text-blue-300 text-sm"
                        >
                          View on Explorer
                        </a>
                      ) : (
                        <span className="text-gray-500 text-sm">Pending</span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {trades && trades.length > 0 && (
        <div className="mt-6 bg-slate-800 rounded-lg p-6 border border-slate-700">
          <h3 className="text-lg font-semibold mb-4">Trade Statistics</h3>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="p-4 bg-slate-700 rounded-lg">
              <div className="text-gray-400 text-sm mb-1">Total Trades</div>
              <div className="text-2xl font-bold">{trades.length}</div>
            </div>
            <div className="p-4 bg-slate-700 rounded-lg">
              <div className="text-gray-400 text-sm mb-1">Buy Orders</div>
              <div className="text-2xl font-bold text-green-400">
                {trades.filter((t) => t.trade_type === 'BUY').length}
              </div>
            </div>
            <div className="p-4 bg-slate-700 rounded-lg">
              <div className="text-gray-400 text-sm mb-1">Sell Orders</div>
              <div className="text-2xl font-bold text-red-400">
                {trades.filter((t) => t.trade_type === 'SELL').length}
              </div>
            </div>
            <div className="p-4 bg-slate-700 rounded-lg">
              <div className="text-gray-400 text-sm mb-1">Total Volume</div>
              <div className="text-2xl font-bold">
                $
                {trades
                  .reduce((sum, trade) => sum + trade.total_amount, 0)
                  .toFixed(2)}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default History
