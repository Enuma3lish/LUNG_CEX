import React from 'react'
import { useQuery } from 'react-query'
import { portfolioService } from '../services/api'
import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip } from 'recharts'

const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899']

const Portfolio = () => {
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

  const pieData = portfolio?.holdings?.map((holding) => ({
    name: holding.asset?.symbol,
    value: holding.value,
  })) || []

  if (portfolio?.cash > 0) {
    pieData.push({
      name: 'Cash',
      value: portfolio.cash,
    })
  }

  const pnlColor = portfolio?.pnl >= 0 ? 'text-green-400' : 'text-red-400'
  const pnlSign = portfolio?.pnl >= 0 ? '+' : ''

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Portfolio</h1>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
          <h3 className="text-gray-400 text-sm font-medium mb-2">Total Value</h3>
          <p className="text-2xl font-bold">${portfolio?.total_value?.toFixed(2)}</p>
        </div>

        <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
          <h3 className="text-gray-400 text-sm font-medium mb-2">Cash</h3>
          <p className="text-2xl font-bold">${portfolio?.cash?.toFixed(2)}</p>
        </div>

        <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
          <h3 className="text-gray-400 text-sm font-medium mb-2">Assets</h3>
          <p className="text-2xl font-bold">
            ${(portfolio?.total_value - portfolio?.cash)?.toFixed(2)}
          </p>
        </div>

        <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
          <h3 className="text-gray-400 text-sm font-medium mb-2">Total P&L</h3>
          <p className={`text-2xl font-bold ${pnlColor}`}>
            {pnlSign}${portfolio?.pnl?.toFixed(2)}
          </p>
          <p className={`text-sm ${pnlColor}`}>
            {pnlSign}
            {((portfolio?.pnl / 10000) * 100).toFixed(2)}%
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
          <h3 className="text-xl font-semibold mb-4">Asset Allocation</h3>
          {pieData.length > 0 ? (
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={pieData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={(entry) => `${entry.name}: ${((entry.value / portfolio.total_value) * 100).toFixed(1)}%`}
                  outerRadius={100}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {pieData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip
                  formatter={(value) => `$${value.toFixed(2)}`}
                  contentStyle={{ backgroundColor: '#1e293b', border: '1px solid #334155' }}
                />
              </PieChart>
            </ResponsiveContainer>
          ) : (
            <div className="text-center text-gray-400 py-12">No assets to display</div>
          )}
        </div>

        <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
          <h3 className="text-xl font-semibold mb-4">Performance Metrics</h3>
          <div className="space-y-4">
            <div className="flex justify-between items-center p-4 bg-slate-700 rounded-lg">
              <span className="text-gray-400">Initial Balance</span>
              <span className="font-semibold">$10,000.00</span>
            </div>
            <div className="flex justify-between items-center p-4 bg-slate-700 rounded-lg">
              <span className="text-gray-400">Current Value</span>
              <span className="font-semibold">${portfolio?.total_value?.toFixed(2)}</span>
            </div>
            <div className="flex justify-between items-center p-4 bg-slate-700 rounded-lg">
              <span className="text-gray-400">Profit/Loss</span>
              <span className={`font-semibold ${pnlColor}`}>
                {pnlSign}${portfolio?.pnl?.toFixed(2)}
              </span>
            </div>
            <div className="flex justify-between items-center p-4 bg-slate-700 rounded-lg">
              <span className="text-gray-400">Return</span>
              <span className={`font-semibold ${pnlColor}`}>
                {pnlSign}
                {((portfolio?.pnl / 10000) * 100).toFixed(2)}%
              </span>
            </div>
            <div className="flex justify-between items-center p-4 bg-slate-700 rounded-lg">
              <span className="text-gray-400">Positions</span>
              <span className="font-semibold">{portfolio?.holdings?.length || 0}</span>
            </div>
          </div>
        </div>
      </div>

      <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
        <h3 className="text-xl font-semibold mb-4">Holdings Detail</h3>
        {portfolio?.holdings?.length === 0 ? (
          <div className="text-center text-gray-400 py-8">No holdings</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="text-left border-b border-slate-700">
                  <th className="pb-3 text-gray-400 font-medium">Asset</th>
                  <th className="pb-3 text-gray-400 font-medium text-right">Quantity</th>
                  <th className="pb-3 text-gray-400 font-medium text-right">Avg Price</th>
                  <th className="pb-3 text-gray-400 font-medium text-right">Current Price</th>
                  <th className="pb-3 text-gray-400 font-medium text-right">Value</th>
                  <th className="pb-3 text-gray-400 font-medium text-right">P&L</th>
                  <th className="pb-3 text-gray-400 font-medium text-right">P&L %</th>
                  <th className="pb-3 text-gray-400 font-medium text-right">Allocation</th>
                </tr>
              </thead>
              <tbody>
                {portfolio?.holdings?.map((holding) => (
                  <tr key={holding.id} className="border-b border-slate-700">
                    <td className="py-4">
                      <div>
                        <div className="font-medium">{holding.asset?.symbol}</div>
                        <div className="text-sm text-gray-400">{holding.asset?.name}</div>
                      </div>
                    </td>
                    <td className="py-4 text-right">{holding.quantity?.toFixed(4)}</td>
                    <td className="py-4 text-right">${holding.avg_price?.toFixed(2)}</td>
                    <td className="py-4 text-right">${holding.current_price?.toFixed(2)}</td>
                    <td className="py-4 text-right font-semibold">
                      ${holding.value?.toFixed(2)}
                    </td>
                    <td
                      className={`py-4 text-right ${
                        holding.pnl >= 0 ? 'text-green-400' : 'text-red-400'
                      }`}
                    >
                      {holding.pnl >= 0 ? '+' : ''}${holding.pnl?.toFixed(2)}
                    </td>
                    <td
                      className={`py-4 text-right ${
                        holding.pnl_percent >= 0 ? 'text-green-400' : 'text-red-400'
                      }`}
                    >
                      {holding.pnl_percent >= 0 ? '+' : ''}
                      {holding.pnl_percent?.toFixed(2)}%
                    </td>
                    <td className="py-4 text-right">
                      {((holding.value / portfolio.total_value) * 100).toFixed(1)}%
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}

export default Portfolio
