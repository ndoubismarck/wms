import Dashboard from '@/views/pages/Admin/Dashboard/Dashboard.vue'
import Staff from '@/views/pages/Admin/Staff/Staff.vue'
import Tasks from '@/views/pages/Admin/Tasks/Tasks.vue'
import AddTask from '@/views/pages/Admin/Tasks/AddTask.vue'
import Teams from '@/views/pages/Admin/Teams/Teams.vue'
import Inventory from '@/views/pages/Admin/Inventory/Inventory.vue'
import InventoryMovements from '@/views/pages/Admin/Inventory/Movements.vue'
import ProductDetails from '@/views/pages/Admin/Inventory/ProductDetails.vue'
import AddProduct from '@/views/pages/Admin/Inventory/AddProduct.vue'
import AddLocationBin from '@/views/pages/Admin/Inventory/AddLocationBin.vue'
import AddLocationAisle from '@/views/pages/Admin/Inventory/AddLocationAisle.vue'
import AddLocationBay from '@/views/pages/Admin/Inventory/AddLocationBay.vue'
import AddLocationShelf from '@/views/pages/Admin/Inventory/AddLocationShelf.vue'
import AddLocationShelfLevel from '@/views/pages/Admin/Inventory/AddLocationShelfLevel.vue'
import AddLocationBinCategory from '@/views/pages/Admin/Inventory/AddLocationBinCategory.vue'
import Products from '@/views/pages/Admin/Inventory/Products.vue'
import Warehouse from '@/views/pages/Admin/Warehouse/Warehouse.vue'
import AddProductBrand from '@/views/pages/Admin/Inventory/AddProductBrand.vue'
import AddProductCategory from '@/views/pages/Admin/Inventory/AddProductCategory.vue'
import AddProductSubcategory from '@/views/pages/Admin/Inventory/AddProductSubcategory.vue'
import Orders from '@/views/pages/Admin/Operations/Orders.vue'
import Shipments from '@/views/pages/Admin/Operations/Shipments.vue'
import Customers from '@/views/pages/Admin/Customers/Customers.vue'
import Settings from '@/views/pages/Admin/Settings/Settings.vue'
import Suppliers from '@/views/pages/Admin/Suppliers/Suppliers.vue'
import Help from '@/views/pages/Admin/Help/Help.vue'

export default [
  {
    path: '/admin',
    name: 'admin:home',
    meta: {
      title: 'Dashboard',
      breadcrumb: [],
    },
    component: Dashboard,
  },
  {
    path: '/admin/inventory',
    name: 'admin:inventory',
    meta: {
      title: 'Inventory',
      breadcrumb: ['admin:inventory'],
    },
    component: Inventory,
  },
  {
    path: '/admin/staff',
    name: 'admin:staff',
    meta: {
      title: 'Staff',
      breadcrumb: ['admin:staff'],
    },
    component: Staff,
  },
  {
    path: '/admin/tasks',
    name: 'admin:tasks',
    meta: {
      title: 'Tasks',
      breadcrumb: ['admin:tasks'],
    },
    component: Tasks,
    children: [
      {
        path: '/admin/tasks/add',
        name: 'admin:tasks:add',
        meta: {
          title: 'Add Task',
          breadcrumb: ['admin:tasks'],
        },
        component: AddTask,
      },
    ],
  },
  {
    path: '/admin/teams',
    name: 'admin:teams',
    meta: {
      title: 'Teams',
      breadcrumb: ['admin:teams'],
    },
    component: Teams,
  },
  {
    path: '/admin/inventory/products',
    name: 'admin:inventory:products',
    meta: {
      title: 'Products',
      breadcrumb: ['admin:inventory', 'admin:inventory:products'],
    },
    component: Products,
    children: [
      {
        path: '/admin/inventory/products/add',
        name: 'admin:inventory:products:add',
        meta: {
          title: 'Add Product',
          breadcrumb: ['admin:inventory'],
        },
        component: AddProduct,
        children: [
          {
            path: '/admin/inventory/products/add/brand',
            name: 'admin:inventory:products:add:brand',
            meta: {
              title: 'Add Product Brand',
              breadcrumb: ['admin:inventory'],
            },
            component: AddProductBrand,
          },
          {
            path: '/admin/inventory/products/add/category/:brandId',
            name: 'admin:inventory:products:add:category',
            meta: {
              title: 'Add Product Category',
              breadcrumb: ['admin:inventory'],
            },
            component: AddProductCategory,
          },
          {
            path: '/admin/inventory/products/add/subcategory/:categoryId',
            name: 'admin:inventory:products:add:subcategory',
            meta: {
              title: 'Add Product Subcategory',
              breadcrumb: ['admin:inventory'],
            },
            component: AddProductSubcategory,
          },
          {
            path: '/admin/inventory/products/add/aisle',
            name: 'admin:inventory:products:aisle:add',
            meta: {
              title: 'Add Aisle',
              breadcrumb: ['admin:inventory'],
            },
            component: AddLocationAisle,
          },
          {
            path: '/admin/inventory/products/add/bay/:aisleId',
            name: 'admin:inventory:products:bay:add',
            meta: {
              title: 'Add Bay',
              breadcrumb: ['admin:inventory'],
            },
            component: AddLocationBay,
          },
          {
            path: '/admin/inventory/products/add/shelf/:bayId',
            name: 'admin:inventory:products:shelf:add',
            meta: {
              title: 'Add Shelf',
              breadcrumb: ['admin:inventory'],
            },
            component: AddLocationShelf,
          },
          {
            path: '/admin/inventory/products/add/shelf/level/:shelfId',
            name: 'admin:inventory:products:shelf:level:add',
            meta: {
              title: 'Add Shelf Level',
              breadcrumb: ['admin:inventory'],
            },
            component: AddLocationShelfLevel,
          },
          {
            path: '/admin/inventory/products/add/bin/category',
            name: 'admin:inventory:products:bin:category:add',
            meta: {
              title: 'Add Bin Category',
              breadcrumb: ['admin:inventory'],
            },
            component: AddLocationBinCategory,
          },
          {
            path: '/admin/inventory/products/add/bin/:shelfLevelId',
            name: 'admin:inventory:products:bin:add',
            meta: {
              title: 'Add Location Bin',
              breadcrumb: ['admin:inventory'],
            },
            component: AddLocationBin,
          },
        ],
      },
      {
        path: '/admin/inventory/products/:id',
        name: 'admin:inventory:products:details',
        meta: {
          title: 'Product',
          breadcrumb: ['admin:inventory'],
        },
        component: ProductDetails,
      },
    ],
  },
  {
    path: '/admin/inventory/products/:id',
    name: 'admin:inventory:products:view',
    meta: {
      title: 'View Product',
      breadcrumb: ['admin:inventory', 'admin:inventory:products:view'],
    },
    component: ProductDetails,
  },
  {
    path: '/admin/inventory/movements',
    name: 'admin:inventory:movements',
    meta: {
      title: 'Movements',
      breadcrumb: ['admin:inventory', 'admin:inventory:movements'],
    },
    component: InventoryMovements,
  },
  {
    path: '/admin/operations/orders',
    name: 'admin:operations:orders',
    meta: {
      title: 'Orders',
      breadcrumb: ['admin:operations:orders'],
    },
    component: Orders,
  },
  {
    path: '/admin/operations/shipments',
    name: 'admin:operations:shipments',
    meta: {
      title: 'Shipments',
      breadcrumb: ['admin:operations:shipments'],
    },
    component: Shipments,
  },
  {
    path: '/admin/operations/customers',
    name: 'admin:suppliers:customers',
    meta: {
      title: 'Customers',
      breadcrumb: ['admin:suppliers:customers'],
    },
    component: Customers,
  },
  {
    path: '/admin/operations/suppliers',
    name: 'admin:suppliers:list',
    meta: {
      title: 'Suppliers',
      breadcrumb: ['admin:suppliers:list'],
    },
    component: Suppliers,
  },
  {
    path: '/admin/warehouse',
    name: 'admin:warehouse',
    meta: {
      title: 'Warehouse',
      breadcrumb: ['admin:warehouse'],
    },
    component: Warehouse,
  },
  {
    path: '/admin/warehouse/bins/:id',
    name: 'admin:warehouse:bins:view',
    meta: {
      title: 'Bins',
      breadcrumb: ['admin:warehouse', 'admin:warehouse:bins'],
    },
    component: Inventory,
  },
  {
    path: '/admin/warehouse/shelves/:id',
    name: 'admin:warehouse:shelves:view',
    meta: {
      title: 'Shelf Levels',
      breadcrumb: ['admin:warehouse', 'admin:warehouse:shelves:view'],
    },
    component: Inventory,
  },
  {
    path: '/admin/warehouse/shelves/levels/:id',
    name: 'admin:warehouse:shelves:levels:view',
    meta: {
      title: 'Shelf Levels',
      breadcrumb: ['admin:warehouse', 'admin:warehouse:shelves:levels:view'],
    },
    component: Inventory,
  },
  {
    path: '/admin/warehouse/bays/:id',
    name: 'admin:warehouse:bays:view',
    meta: {
      title: 'Bays',
      breadcrumb: ['admin:warehouse', 'admin:warehouse:bays'],
    },
    component: Inventory,
  },
  {
    path: '/admin/warehouse/aisles/:id',
    name: 'admin:warehouse:aisles:view',
    meta: {
      title: 'Aisles',
      breadcrumb: ['admin:warehouse', 'admin:warehouse:aisles'],
    },
    component: Inventory,
  },
  {
    path: '/admin/account',
    name: 'admin:account',
    meta: {
      title: 'Account',
      breadcrumb: ['admin:account'],
    },
    component: Dashboard,
  },
  {
    path: '/admin/settings',
    name: 'admin:settings',
    meta: {
      title: 'Settings',
      breadcrumb: ['admin:settings'],
    },
    component: Settings,
  },
  {
    path: '/admin/help',
    name: 'admin:help',
    meta: {
      title: 'Help and Support',
    },
    component: Help,
  },
  {
    path: '/admin/activation',
    name: 'admin:activation',
    meta: {
      title: 'Product Activation',
    },
    component: Dashboard,
  },
]
