import React from 'react'
import { useQuery } from 'react-query'
import { useNavigate } from 'react-router-dom'
import { portfolioService } from '../services/api'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

const Dashboard = () => {
  const navigate = useNavigate()
  const { data: portfolio, isLoading } = useQuery('portfolio', portfolioService.getPortfolio, {
    refetchInterval: 5000,
  })

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center text-gray-400">Loading...</div>
      </div>
    )
  }

  const pnlColor = portfolio?.pnl >= 0 ? 'text-green-400' : 'text-red-400'
  const pnlSign = portfolio?.pnl >= 0 ? '+' : ''

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
          <h3 className="text-gray-400 text-sm font-medium mb-2">Total Portfolio Value</h3>
          <p className="text-3xl font-bold">${portfolio?.total_value?.toFixed(2)}</p>
          <p className={`text-sm mt-2 ${pnlColor}`}>
            {pnlSign}${portfolio?.pnl?.toFixed(2)} ({pnlSign}
            {((portfolio?.pnl / 10000) * 100).toFixed(2)}%)
          </p>
        </div>

        <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
          <h3 className="text-gray-400 text-sm font-medium mb-2">Cash Balance</h3>
          <p className="text-3xl font-bold">${portfolio?.cash?.toFixed(2)}</p>
          <p className="text-sm text-gray-400 mt-2">Available for trading</p>
        </div>

        <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
          <h3 className="text-gray-400 text-sm font-medium mb-2">Assets Value</h3>
          <p className="text-3xl font-bold">
            ${(portfolio?.total_value - portfolio?.cash)?.toFixed(2)}
          </p>
          <p className="text-sm text-gray-400 mt-2">{portfolio?.holdings?.length || 0} positions</p>
        </div>
      </div>

      <div className="bg-slate-800 rounded-lg p-6 border border-slate-700 mb-8">
        <h3 className="text-xl font-semibold mb-4">Holdings</h3>
        {portfolio?.holdings?.length === 0 ? (
          <div className="text-center py-8">
            <p className="text-gray-400 mb-4">You don't have any holdings yet</p>
            <button
              onClick={() => navigate('/trading')}
              className="px-6 py-2 bg-blue-600 hover:bg-blue-700 rounded-md font-medium"
            >
              Start Trading
            </button>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="text-left border-b border-slate-700">
                  <th className="pb-3 text-gray-400 font-medium">Asset</th>
                  <th className="pb-3 text-gray-400 font-medium">Quantity</th>
                  <th className="pb-3 text-gray-400 font-medium">Avg Price</th>
                  <th className="pb-3 text-gray-400 font-medium">Current Price</th>
                  <th className="pb-3 text-gray-400 font-medium">Value</th>
                  <th className="pb-3 text-gray-400 font-medium">P&L</th>
                </tr>
              </thead>
              <tbody>
                {portfolio?.holdings?.map((holding) => (
                  <tr key={holding.id} className="border-b border-slate-700">
                    <td className="py-4 font-medium">{holding.asset?.symbol}</td>
                    <td className="py-4">{holding.quantity?.toFixed(4)}</td>
                    <td className="py-4">${holding.avg_price?.toFixed(2)}</td>
                    <td className="py-4">${holding.current_price?.toFixed(2)}</td>
                    <td className="py-4">${holding.value?.toFixed(2)}</td>
                    <td className={`py-4 ${holding.pnl >= 0 ? 'text-green-400' : 'text-red-400'}`}>
                      {holding.pnl >= 0 ? '+' : ''}${holding.pnl?.toFixed(2)} (
                      {holding.pnl_percent?.toFixed(2)}%)
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
          <h3 className="text-xl font-semibold mb-4">Quick Actions</h3>
          <div className="space-y-3">
            <button
              onClick={() => navigate('/trading')}
              className="w-full px-4 py-3 bg-blue-600 hover:bg-blue-700 rounded-md font-medium text-left"
            >
              Trade Assets
            </button>
            <button
              onClick={() => navigate('/portfolio')}
              className="w-full px-4 py-3 bg-slate-700 hover:bg-slate-600 rounded-md font-medium text-left"
            >
              View Portfolio
            </button>
            <button
              onClick={() => navigate('/history')}
              className="w-full px-4 py-3 bg-slate-700 hover:bg-slate-600 rounded-md font-medium text-left"
            >
              Trade History
            </button>
          </div>
        </div>

        <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
          <h3 className="text-xl font-semibold mb-4">Available Assets</h3>
          <div className="space-y-2">
            <div className="flex justify-between items-center py-2 border-b border-slate-700">
              <span className="font-medium">BTC</span>
              <span className="text-green-400">$45,000</span>
            </div>
            <div className="flex justify-between items-center py-2 border-b border-slate-700">
              <span className="font-medium">ETH</span>
              <span className="text-green-400">$2,500</span>
            </div>
            <div className="flex justify-between items-center py-2 border-b border-slate-700">
              <span className="font-medium">SOL</span>
              <span className="text-green-400">$100</span>
            </div>
            <div className="flex justify-between items-center py-2 border-b border-slate-700">
              <span className="font-medium">USDC</span>
              <span className="text-gray-400">$1.00</span>
            </div>
            <div className="flex justify-between items-center py-2">
              <span className="font-medium">USDT</span>
              <span className="text-gray-400">$1.00</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Dashboard
