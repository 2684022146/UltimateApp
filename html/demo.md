## 一、核心说明

后端仅返回JSON原始数据，前端负责页面渲染

## 二、GEO优化整体规划

### 2.1 优化核心目标

1. 零售端：提升用户搜索时的页面匹配度与排名，强化区域服务覆盖优势。
2. 工程端：触达装修公司、地产商，凸显区域供货与安装能力。
3. 全站：通过规范地理信息呈现，强化定位认知，积累地理相关搜索权重，扩大客群覆盖。
4. 体验&权威度：FAQ、权威资质标注，提升用户信任度与搜索引擎对网站专业度、权威性的识别。

### 2.2 核心优化方向

1. 服务区域优化：明确标注服务覆盖范围，在各产品页关联对应服务区域。
2. 本地服务优化：清晰呈现区域内上门测量、安装、售后等服务时效，强化“本地服务”属性，区别于跨区域电商。
3. 工程对接优化：针对工程类客户，标注工程案例、本地对接团队及服务范围，适配批量采购与定制需求。
4. 页面地理信息可视化：在首页、联系页、产品分类页展示服务覆盖地图、本地案例，提升用户信任度与搜索引擎识别度。
5. 权威信息强化：FAQ结构化数据解答用户高频地理/服务问题，企业权威资质、行业认证标注，提升网站权威性。

### 2.3 前后端分工边界

- 后端：提供地理相关原始数据，包括但不限于核心运营地址、服务覆盖清单、区域服务时效、本地对接联系方式、工程案例所在区域等，确保数据精准、格式统一；提供FAQ问答库、企业权威资质（认证名称、颁发机构、有效期）等原始数据。
- 前端：对接后端数据，在指定页面展示地理信息；按规范构建JSON-LD结构化数据（含原有地理字段+FAQ、权威信息字段），渲染至页面头部（`<head>`标签内），确保搜索引擎可抓取；优化页面地理信息呈现形式，提升用户体验。

## 三、JSON-LD 结构化数据实施规范

JSON-LD是搜索引擎识别网站地理属性、权威属性的核心载体，需结合零售、批发不同场景，选择对应数据类型，填充完整字段，确保与页面内容一致。前端需将JSON-LD以`<script>`标签形式渲染至页面头部（`<head>`标签内），标签类型固定为“application/ld+json”，禁止放入异步加载文件，避免抓取失败。标签核心格式如下：

**基础标签格式**：`<script type="application/ld+json">{JSON-LD数据内容}</script>`，其中JSON-LD数据需严格遵循JSON格式规范，字段值与页面内容、后端数据完全对齐。

### 3.1 核心适用JSON-LD类型及场景

#### 3.1.1 商品地理属性标注（核心场景）

适用页面：橱柜分类页、地板页、声学板页、五金配件页等所有产品页面，核心类型为“Product”“CollectionPage”，嵌套地理相关字段、FAQ字段、权威信息字段，标注服务覆盖区域、本地服务属性、用户高频问题、企业权威资质，强化商品与地理区域的关联及网站权威性。

**首页Organization+WebSite+FAQ+权威信息组合Schema**：

```html

<script type="application/ld+json">
{
    "@context": "https://schema.org",
    "@graph": [
        {
            "@type": "Organization",//企业信息核心节点
            "@id": "https://www.everlastingcabinetry.com/#organization",//全网唯一标识
            "name": "Everlasting Cabinetry",//企业官方名称
            "url": "https://www.everlastingcabinetry.com/",//企业官网地址
            "telephone": "+1-402-932-3932",//企业联系电话
            "description": "Everlasting Cabinetry provides premium custom cabinetry for homes and projects across the Midwest, offering high-quality cabinets, flooring, hardware, and more to elevate any space.",//企业核心描述
            "address": {//企业物理地址
                "@type": "PostalAddress",//邮政地址类型
                "streetAddress": "10025 I St",
                "addressLocality": "Omaha",
                "addressRegion": "NE",
                "postalCode": "68127",
                "addressCountry": "US"
            },
            "areaServed": [//企业服务覆盖区域
                {"@type": "Place", "name": "Colorado"},
                {"@type": "Place", "name": "Iowa"},
                {"@type": "Place", "name": "Kansas"},
                {"@type": "Place", "name": "Missouri"},
                {"@type": "Place", "name": "Nebraska"},
                {"@type": "Place", "name": "North Dakota"},
                {"@type": "Place", "name": "South Dakota"}
            ],
            "openingHoursSpecification": [  //营业时间
                {
                    "@type": "OpeningHoursSpecification",
                    "dayOfWeek": ["Monday","Tuesday","Wednesday","Thursday","Friday"],
                    "opens": "09:00",
                    "closes": "17:30"
                },
                {
                    "@type": "OpeningHoursSpecification",
                    "dayOfWeek": ["Saturday","Sunday"],
                    "opens": "00:00",
                    "closes": "00:00"
                }
            ],
            "sameAs": ["https://www.instagram.com/everlastingcabinetryllc"],//企业的社交媒体账号
            "hasCredential": [//：企业权威资质/认证（提升网站权威性）
                {
                    "@type": "OrganizationRole",//资质关联类型
                    "credential": {
                        "@type": "Certificate",//认证类型：证书
                        "name": "National Kitchen & Bath Association (NKBA) Certification",//认证名称：全美厨卫协会认证
                        "issuer": "National Kitchen & Bath Association",//颁发机构
                        "validFrom": "2020-01-15",//认证生效时间
                        "validUntil": "2025-01-14"//认证失效时间
                    },
                    "roleName": "Certified Cabinetry Provider"//资质对应的企业角色
                },
                {
                    "@type": "OrganizationRole",
                    "credential": {
                        "@type": "Award",//认证类型：行业奖项
                        "name": "Midwest Best Custom Cabinetry Provider 2023",//奖项名称：2023中西部最佳定制橱柜供应商
                        "issuer": "Midwest Home & Construction Magazine",//颁发机构
                        "validFrom": "2023-03-01"//奖项生效时间（无失效期则仅填validFrom）
                    },
                    "roleName": "Award-Winning Cabinetry Brand"
                }
            ]
        },
        {
            "@type": "WebSite",//官网的基础属性
            "@id": "https://www.everlastingcabinetry.com/#website",//网站的唯一标识符
            "name": "Everlasting Cabinetry",
            "url": "https://www.everlastingcabinetry.com/",
            "publisher": {"@id": "https://www.everlastingcabinetry.com/#organization"},//网站的发布者（关联上方的 Organization 对象）
            "inLanguage": "en-US",//网站的主要语言
            "potentialAction": {//声明网站有「搜索功能」
                "@type": "SearchAction",
                "target": "https://www.everlastingcabinetry.com/?s={search_term_string}",
                "query-input": "required name=search_term_string"
            }
        },
        {
            "@type": "HomePage",//首页
            "@id": "https://www.everlastingcabinetry.com/#homepage",//首页的唯一标识符
            "url": "https://www.everlastingcabinetry.com/",
            "name": "Everlasting Cabinetry – Premium Custom Cabinetry Across the Midwest",//对应title标签
            "isPartOf": {"@id": "https://www.everlastingcabinetry.com/#website"},//该页面所属的网站
            "about": {"@id": "https://www.everlastingcabinetry.com/#organization"},//该页面的核心主题
            "inLanguage": "en-US",
            "description": "Everlasting Cabinetry provides premium custom cabinetry for homes and projects across the Midwest. Explore high-quality cabinets, flooring, hardware, and more to elevate any space.",//页面的描述
            "mainEntity": {//：首页核心FAQ模块（解答用户高频地理/服务问题）
                "@type": "FAQPage",//类型：FAQ页面
                "mainEntity": [//FAQ列表，支持多个问答对
                    {
                        "@type": "Question",//问答对-问题
                        "name": "What states do you offer cabinet delivery and installation services?",//问题内容（用户高频问的地理服务范围问题）
                        "acceptedAnswer": {
                            "@type": "Answer",//问答对-答案
                            "text": "We provide delivery and installation services across 7 Midwest states: Colorado, Iowa, Kansas, Missouri, Nebraska, North Dakota, and South Dakota. Nebraska customers enjoy faster 2-3 business day delivery, while other states receive deliveries within 5-7 business days."//答案内容，需包含地理关键词，与页面展示的FAQ一致
                        }
                    },
                    {
                        "@type": "Question",
                        "name": "Do you offer wholesale pricing for commercial cabinetry projects?",//批发场景高频问题
                        "acceptedAnswer": {
                            "@type": "Answer",
                            "text": "Yes! We offer competitive wholesale pricing for commercial and engineering-grade cabinetry projects across the Midwest. The minimum order quantity for wholesale orders is 50 units, and we provide on-site measurement and custom installation services in Colorado, Nebraska, and Missouri."
                        }
                    },
                    {
                        "@type": "Question",
                        "name": "How long does custom cabinet installation take in Omaha, NE?",//本地服务时效高频问题
                        "acceptedAnswer": {
                            "@type": "Answer",
                            "text": "For custom cabinet installation in Omaha, NE, our team typically completes the project within 3-5 business days after delivery, depending on the size and complexity of the project."
                        }
                    }
                ]
            }
        }
    ]
}
</script>
```

**有框橱柜页Schema：集合页（FAQ+权威信息关联）**：

```html

<script type="application/ld+json">
{
    "@context": "https://schema.org", // 声明结构化数据遵循的上下文标准（固定为schema.org）
    "@type": "CollectionPage", // 页面类型：集合页（对应橱柜分类列表页，区别于单一商品详情页）
    "@id": "https://www.everlastingcabinetry.com/category?category_cd=ELC0001FC#collection", // 该页面的唯一标识ID，用于层级关联
    "url": "https://www.everlastingcabinetry.com/category?category_cd=ELC0001FC", // 页面的完整URL地址
    "name": "Framed Cabinetry", // 页面名称（对应H1核心关键词，需与页面标题一致）
    "inLanguage": "en-US", // 页面使用的语言（英文，美国地区）
    "description": "Explore premium framed cabinetry designed for kitchens, bathrooms, and custom projects. Shop high-quality framed cabinet styles from Everlasting Cabinetry.", // 页面描述融入产品核心卖点
    "isPartOf": {"@id": "https://www.everlastingcabinetry.com/category?category_cd=ELC0001#collection"}, // 关联上级页面ID（表示该有框橱柜页属于橱柜总分类页）
    "about": {"@id": "https://www.everlastingcabinetry.com/#organization"}, // 关联企业主体ID（强化页面与品牌的归属关系，继承企业权威资质）
    "mainEntity": { // 页面核心实体信息（该分类页的商品目录核心内容+FAQ）
        "@type": "ItemList",// 整合商品目录+FAQ的列表类型
        "itemListElement": [
            {
                "@type": "ListItem", // 列表项1：商品目录信息
                "position": 1, // 排序位置
                "item": {
                    "@type": "OfferCatalog", // 实体类型：商品目录（适配分类页批量商品展示场景）
                    "name": "Framed Cabinetry Products", // 商品目录名称（对应分类页核心商品集合名称）
                    "itemListElement": [ // 目录下的商品列表项（支持多个ListItem，此处为分类级核心项）
                        {
                            "@type": "ListItem", // 列表项类型（标识为目录中的单个列表条目）
                            "position": 1, // 列表项排序位置（多个条目时按顺序标注，此处仅1个核心分类项）
                            "itemOffered": { // 列表项提供的商品/服务核心信息
                                "@type": "Product", // 商品类型（标识为产品，区别于服务类）
                                "name": "Framed Cabinetry", // 商品名称（分类级商品名称，非单一SKU）
                                "category": "Cabinetry", // 商品所属大品类（用于搜索引擎识别商品类目）
                                "areaServed": [ // GEO优化核心字段：商品服务覆盖区域（适配地理相关搜索）
                                    {"@type": "Place", "name": "Colorado"}, // 覆盖州1：科罗拉多州
                                    {"@type": "Place", "name": "Nebraska"}, // 覆盖州2：内布拉斯加州
                                    {"@type": "Place", "name": "Missouri"} // 覆盖州3：密苏里州
                                ],
                                "hasMerchantReturnPolicy": {//可选：退货政策（强化本地服务属性）
                                    "@type": "MerchantReturnPolicy",
                                    "applicableRegion": ["Colorado","Nebraska","Missouri"],//退货政策适用区域
                                    "returnMethod": "https://schema.org/ReturnInPerson",//退货方式：到店退货
                                    "returnPeriod": {
                                        "@type": "QuantitativeValue",
                                        "value": 30,//退货周期30天
                                        "unitCode": "DAY"
                                    }
                                }
                            }
                        }
                    ]
                }
            },
            {
                "@type": "ListItem", // 列表项2：分类页专属FAQ
                "position": 2,
                "item": {
                    "@type": "FAQPage",
                    "mainEntity": [
                        {
                            "@type": "Question",
                            "name": "What materials are your framed cabinets made of?",//分类页商品相关问题
                            "acceptedAnswer": {
                                "@type": "Answer",
                                "text": "Our framed cabinets are crafted from high-quality solid wood and plywood, sourced from sustainable forests in the Midwest. All materials meet NKBA certification standards for durability and safety."//答案关联权威认证（NKBA），提升可信度
                            }
                        },
                        {
                            "@type": "Question",
                            "name": "Do you offer custom sizing for framed cabinets in Iowa?",//地理+商品定制问题
                            "acceptedAnswer": {
                                "@type": "Answer",
                                "text": "Yes, we offer fully custom sizing for framed cabinets in Iowa and all our service states. Our team can create cabinets tailored to your kitchen or bathroom dimensions, with free on-site measurements available for Iowa customers."
                            }
                        }
                    ]
                }
            }
        ]
    }
}
</script>
```

**零售商品完整标签示例（FAQ+权威信息）**：商品详情页
```html
<script type="application/ld+json">
{
    "@context": "https://schema.org", // 声明结构化数据遵循的核心上下文标准（固定为schema.org）
    "@type": "Product", // 页面/实体类型：零售商品（区别于集合页、服务类等类型）
    "name": "XX品牌运动鞋（男款）", // 商品名称（需与页面展示的商品名称完全一致，含品牌/品类/规格）
    "sku": "SP-M-2024001", // 商品唯一编码（SKU，零售商品必选字段，用于标识具体商品规格）
    "availability": "https://schema.org/InStock", // 库存状态（枚举值，InStock=有货，OutOfStock=缺货，需与页面库存展示一致）
    "availableAtOrFrom": { // 核心GEO字段：商品发货/可获取的地点（适配零售场景地理属性标注）
        "@type": "Place", // 子类型：地点（用于定义发货仓/自提点等物理位置）
        "name": "XX网站北京通州仓", // 发货仓/自提点名称（便于搜索引擎识别具体仓储节点）
        "address": { // 地点的详细邮政地址（GEO优化核心，需精准到区/县级别）
            "@type": "PostalAddress", // 子类型：邮政地址（遵循schema.org地址规范）
            "streetAddress": "北京市通州区XX物流园A区3栋", // 街道级详细地址（需完整、准确）
            "addressLocality": "通州区", // 地址所属区县（GEO精准化关键字段）
            "addressRegion": "北京市", // 地址所属省市/州（适配区域搜索匹配）
            "postalCode": "101100", // 邮政编码（提升地址精准度，辅助搜索引擎地理识别）
            "addressCountry": "CN" // 国家代码（ISO 3166-1 alpha-2格式，CN=中国，US=美国等）
        },
        "geo": { // 地理坐标（可选但推荐，进一步强化GEO属性，适配地图类搜索）
            "@type": "GeoCoordinates", // 子类型：地理坐标
            "latitude": 39.928800, // 纬度（精准至6位小数，需与发货仓实际位置一致）
            "longitude": 116.657200 // 经度（精准至6位小数，需与发货仓实际位置一致）
        }
    },
    "offers": { // 商品报价信息（零售场景核心交易属性）
        "@type": "Offer", // 子类型：报价
        "price": "399.00", // 商品售价（需与页面展示价格一致，保留2位小数）
        "priceCurrency": "CNY" // 价格货币类型（ISO 4217代码，CNY=人民币，USD=美元等）
    },
    "certification": [//：商品权威认证（提升商品可信度）
        {
            "@type": "Certificate",
            "name": "国家质量认证中心产品质量认证",//认证名称
            "issuer": "中国国家质量认证中心",//颁发机构
            "validFrom": "2023-01-01"//认证生效时间
        }
    ],
    "mainEntity": {//：商品详情页FAQ（解答用户高频问题）
        "@type": "FAQPage",
        "mainEntity": [
            {
                "@type": "Question",
                "name": "这款运动鞋从北京通州仓发货后，多久能到朝阳区？",//地理+配送时效问题
                "acceptedAnswer": {
                    "@type": "Answer",
                    "text": "从北京通州仓发货至朝阳区，同城配送时效为1个工作日，下单后当天16点前付款可次日送达。"
                }
            },
            {
                "@type": "Question",
                "name": "商品是否支持北京地区到店自提？",//本地服务问题
                "acceptedAnswer": {
                    "@type": "Answer",
                    "text": "支持的，您可选择北京通州仓自提点到店自提，自提时间为周一至周五9:00-17:30，周末暂不支持自提。"
                }
            }
        ]
    }
}
</script>
```

**批发商品完整标签示例（FAQ+权威信息）**：
```html
<script type="application/ld+json">
{
    "@context": "https://schema.org", // 声明结构化数据遵循的核心上下文标准（固定为schema.org，搜索引擎通用解析规范）
    "@type": "Product", // 实体类型：单一商品（批发场景下核心标注类型，区别于集合页/服务类）
    "name": "XX品牌建筑防水涂料（50kg/桶）", // 商品名称（需与页面展示名称完全一致，含品牌、品类、规格，适配批发场景规格标注）
    "sku": "PF-B-2024005", // 商品唯一编码（批发商品必选字段，PF前缀标识批发类型，便于区分零售SKU）
    "availability": "https://schema.org/InStock", // 库存状态（枚举值，InStock=有货，需与对应发货仓库存状态一致，批发场景需关联区域库存）
    "availableAtOrFrom": { // 核心GEO字段：批发商品发货/可获取的物理地点（适配GEO优化，精准标注发货仓）
        "@type": "Place", // 子类型：地点（定义发货仓的物理位置属性）
        "name": "XX网站天津武清仓", // 发货仓名称（便于搜索引擎识别具体仓储节点，批发场景需明确仓点名称）
        "address": { // 发货仓详细邮政地址（GEO优化核心，需精准到区/县级别）
            "@type": "PostalAddress", // 子类型：邮政地址（遵循schema.org地址规范）
            "streetAddress": "天津市武清区XX工业园B区8号", // 街道级详细地址（需完整、准确，批发场景便于物流对接）
            "addressLocality": "武清区", // 地址所属区县（GEO精准化关键字段，批发配送范围匹配核心）
            "addressRegion": "天津市", // 地址所属省市（适配区域批发覆盖范围识别）
            "postalCode": "301700", // 邮政编码（提升地址精准度，辅助搜索引擎地理识别）
            "addressCountry": "CN" // 国家代码（ISO 3166-1 alpha-2格式，CN=中国，批发场景明确跨境/境内属性）
        },
        "geo": { // 地理坐标（批发场景推荐标注，强化GEO属性，适配地图类批发商户搜索）
            "@type": "GeoCoordinates", // 子类型：地理坐标
            "latitude": 39.376400, // 纬度（精准至6位小数，需与发货仓实际位置一致）
            "longitude": 117.049700 // 经度（精准至6位小数，需与发货仓实际位置一致）
        }
    },
    "offers": { // 批发商品核心报价信息（批发场景专属嵌套字段，标注批量采购规则）
        "@type": "Offer", // 子类型：报价（适配批发场景的交易属性标注）
        "price": "280.00", // 批发单价（需与页面展示的批量采购单价一致，保留2位小数）
        "priceCurrency": "CNY", // 价格货币类型（ISO 4217代码，CNY=人民币）
        "eligibleRegion": "京津冀地区", // 批发覆盖区域（批发场景必选，明确可供货的地理范围，适配“京津冀防水涂料批发”等搜索）
        "minimumOrderQuantity": 10, // 批发最小起订量（批发场景核心字段，标注最低采购数量，区分零售场景）
        "availableDeliveryMethod": "https://schema.org/LocalDelivery" // 配送方式（枚举值，LocalDelivery=本地配送，适配批发场景的配送服务标注）
    },
    "certification": [//：批发商品权威认证（工程/批发场景核心，提升企业信任度）
        {
            "@type": "Certificate",
            "name": "建筑防水材料国家检测中心合格认证",//认证名称
            "issuer": "国家建筑材料测试中心",//颁发机构
            "validFrom": "2022-06-01",
            "validUntil": "2027-05-31"
        }
    ],
    "mainEntity": {//：批发商品FAQ（工程客户高频问题）
        "@type": "FAQPage",
        "mainEntity": [
            {
                "@type": "Question",
                "name": "京津冀地区批发订货后，多久能配送到位？",//批发配送时效问题
                "acceptedAnswer": {
                    "@type": "Answer",
                    "text": "京津冀地区批发订货，天津武清仓发货后，河北省内1-2个工作日送达，北京、天津同城次日达，最小起订量10桶及以上可免费配送至指定工地。"
                }
            },
            {
                "@type": "Question",
                "name": "工程采购是否支持先验货再付款？",//工程采购信任度问题
                "acceptedAnswer": {
                    "@type": "Answer",
                    "text": "支持的，针对京津冀地区工程采购客户，订单金额满5万元可享受先验货再付款服务，验货地点可选择我方天津武清仓或贵司指定工地。"
                }
            }
        ]
    }
}
</script>
```

必选字段：商品名称、SKU（或商品编码）、库存状态、发货地（含详细地址、区域）、地理坐标（经纬度）；必选（权威/FAQ）：核心页面需包含FAQ模块（至少2组问答）、企业/商品权威资质（至少1项）；可选字段：产地、区域定价、本地退货政策、批发最小起订量（批发商品专属）、商品认证信息。
核心说明：发货地需精准到区/县级别，经纬度精准至6位小数；库存状态需与对应区域关联（如“北京仓现货”“上海仓缺货”），批发商品需额外标注覆盖的批发区域；FAQ需包含地理相关问题，答案需与页面展示内容一致，权威资质需真实有效，禁止虚构。

#### 3.1.2 配送服务标注（FAQ+权威信息，适配配送场景）
适用页面：首页、商品列表页、商品详情页、配送说明页，核心类型为“DeliveryService”，标注配送范围、上门服务时效、安装服务覆盖区域、配送相关FAQ、配送服务资质，适配“本地橱柜配送”“区域安装服务”等搜索需求。

**完整标签示例（适配Everlasting Cabinetry业务，FAQ+权威）**：
```html
<script type="application/ld+json">
{
    "@context": "https://schema.org", // 声明结构化数据遵循的核心标准（schema.org为搜索引擎通用解析规范）
    "@type": "DeliveryService", // 实体类型：配送服务（适配"本地橱柜配送""区域安装服务"等地理相关搜索）
    "name": "Everlasting Cabinetry Midwest Delivery & Installation Service", // 配送服务名称（需包含品牌+覆盖区域+服务类型，强化地理关联）
    "areaServed": [ // 核心GEO字段：配送/安装服务覆盖的地理区域（适配Midwest 7个州的业务范围）
        {"@type": "Place", "name": "Colorado"}, // 子类型：地理地点，标注覆盖的科罗拉多州
        {"@type": "Place", "name": "Iowa"}, // 爱荷华州
        {"@type": "Place", "name": "Kansas"}, // 堪萨斯州
        {"@type": "Place", "name": "Missouri"}, // 密苏里州
        {"@type": "Place", "name": "Nebraska"}, // 内布拉斯加州（核心服务州，配送时效更优）
        {"@type": "Place", "name": "North Dakota"}, // 北达科他州
        {"@type": "Place", "name": "South Dakota"} // 南达科他州
    ],
    "deliveryTime": { // 配送时效规则（区分核心州与其他州，强化本地服务属性）
        "@type": "OpeningHoursSpecification", // 子类型：营业时间/时效规范（复用该类型标注配送时效）
        "dayOfWeek": ["Monday","Tuesday","Wednesday","Thursday","Friday"], // 配送服务可执行的工作日
        "opens": "09:00", // 每日配送服务开始时间
        "closes": "17:30", // 每日配送服务结束时间
        "description": "Delivery within Nebraska: 2-3 business days; Delivery to other Midwest states: 5-7 business days" // 差异化时效说明（核心州更快，适配本地搜索需求）
    },
    "availableDeliveryMethod": "https://schema.org/LocalDelivery", // 配送方式：本地配送（枚举值，区别于快递/自提，强化本地服务标签）
    "hasOfferCatalog": { // 嵌套字段：配送服务包含的增值服务目录（此处为橱柜安装服务）
        "@type": "OfferCatalog", // 子类型：服务目录
        "name": "Midwest Installation Services", // 安装服务目录名称（关联Midwest区域）
        "itemListElement": [ // 服务列表项（支持多个服务，按position排序）
            {
                "@type": "ListItem", // 子类型：列表项
                "position": 1, // 列表排序（多个服务时按优先级标注）
                "itemOffered": { // 具体提供的服务
                    "@type": "Service", // 子类型：服务（区别于商品）
                    "name": "Cabinet Installation Service", // 服务名称：橱柜安装服务
                    "areaServed": ["Colorado","Iowa","Nebraska","Missouri"], // 安装服务覆盖的核心州（精准标注，适配"XX州橱柜安装"搜索）
                    "certification": [//：安装服务权威资质
                        {
                            "@type": "Certificate",
                            "name": "Midwest Contractor License for Cabinet Installation",//安装资质名称
                            "issuer": "Nebraska Department of Labor",//颁发机构（内布拉斯加州劳工部）
                            "validFrom": "2021-05-01"
                        }
                    ]
                }
            }
        ]
    },
    "mainEntity": {//：配送服务FAQ
        "@type": "FAQPage",
        "mainEntity": [
            {
                "@type": "Question",
                "name": "Do you offer emergency delivery for cabinetry in Omaha, NE?",//本地应急配送问题
                "acceptedAnswer": {
                    "@type": "Answer",
                    "text": "Yes, we offer emergency delivery services for cabinetry in Omaha, NE for urgent projects. Emergency deliveries are available Monday-Friday with a 24-hour notice, and a small premium fee applies."
                }
            },
            {
                "@type": "Question",
                "name": "Is installation included with delivery in Colorado?",//配送+安装关联问题
                "acceptedAnswer": {
                    "@type": "Answer",
                    "text": "Installation is included with delivery for all cabinetry orders over $5,000 in Colorado. For smaller orders, we offer discounted installation rates starting at $150 per project."
                }
            }
        ]
    }
}
</script>
```

#### 3.1.3 批发业务专属标注（FAQ+权威信息，批发场景）
适用页面：批发专区首页、工程合作页、橱柜分类页（工程款），核心类型为“Offer”（嵌套于Product类型中），标注批发覆盖范围、工程对接方式、本地测量安装服务、批发FAQ、工程资质，适配工程客户需求。

**工程批发专属标签示例（FAQ+权威）**：
```html
<script type="application/ld+json">
{
    "@context": "https://schema.org", // 声明结构化数据遵循的核心标准（schema.org为搜索引擎通用解析规范，工程批发场景需严格遵循）
    "@type": "Product", // 实体类型：工程批发类商品（适配"工程橱柜批发"等精准搜索，区别于零售商品标注）
    "name": "Commercial Framed Cabinetry (Engineering Grade)", // 商品名称（需包含「商用/工程级」标识，适配工程客户搜索习惯，如"Commercial Cabinetry wholesale"）
    "sku": "ELC-ENG-FC-001", // 工程商品专属SKU（ENG前缀标识工程级，FC=有框橱柜，编码规则需与后端工程商品库一致）
    "offers": { // 工程批发核心报价/服务规则（批发场景核心嵌套字段，标注批量采购条件与配套服务）
        "@type": "Offer", // 子类型：批发报价（适配工程客户批量采购的交易属性标注）
        "price": "220.00", // 工程批发单价（需与工程报价单、页面展示的批量采购单价一致，保留2位小数，货币为USD）
        "priceCurrency": "USD", // 价格货币类型（ISO 4217代码，USD=美元，适配美国中西部工程业务场景）
        "eligibleRegion": ["Colorado","Iowa","Kansas","Missouri","Nebraska","North Dakota","South Dakota"], // 批发供货覆盖区域（工程场景必选，明确7个核心服务州，适配"Midwest commercial cabinetry wholesale"搜索）
        "minimumOrderQuantity": 50, // 批发最小起订量（工程场景核心字段，标注最低采购数量，区别于零售场景，需与业务规则一致）
        "availableDeliveryMethod": "https://schema.org/LocalDelivery", // 配送方式：本地配送（枚举值，适配工程场景的区域配送属性，强化"本地工程供货"标签）
        "includesService": { // 工程批发配套增值服务（工程场景专属字段，标注含上门测量/定制安装，提升工程客户转化）
            "@type": "Service", // 子类型：服务（定义工程配套服务属性）
            "name": "On-Site Measurement & Custom Installation", // 服务名称：上门测量+定制安装（工程橱柜核心服务，适配"Omaha cabinet installation"等地理搜索）
            "areaServed": ["Colorado","Nebraska","Missouri"], // 配套服务覆盖州（精准标注可提供上门服务的核心州，适配"Colorado on-site cabinet measurement"搜索）
            "certification": [//：工程服务权威资质
                {
                    "@type": "Certificate",
                    "name": "Commercial Construction Service Certification",//商用建筑服务认证
                    "issuer": "Midwest Construction Association",//颁发机构
                    "validFrom": "2022-01-01",
                    "validUntil": "2026-12-31"
                }
            ]
        }
    },
    "areaServed": [ // 商品核心服务覆盖区域（GEO优化核心字段，标注商品主要触达的工程客户区域）
        {"@type": "Place", "name": "Colorado"}, // 子类型：地理地点，标注科罗拉多州（工程业务核心覆盖州）
        {"@type": "Place", "name": "Iowa"}, // 爱荷华州（工程批发覆盖州，适配本地装修公司/地产商搜索）
        {"@type": "Place", "name": "Nebraska"} // 内布拉斯加州（核心运营州，强化"Omaha commercial cabinetry"等本地工程搜索匹配度）
    ],
    "mainEntity": {//：工程批发FAQ
        "@type": "FAQPage",
        "mainEntity": [
            {
                "@type": "Question",
                "name": "What payment terms are available for commercial cabinetry wholesale in Missouri?",//工程付款条款问题
                "acceptedAnswer": {
                    "@type": "Answer",
                    "text": "For commercial cabinetry wholesale orders in Missouri, we offer net-30 payment terms for established contractors and developers. A 50% deposit is required upfront, with the remaining balance due upon delivery."
                }
            },
            {
                "@type": "Question",
                "name": "Can you provide custom engineering-grade cabinets for large-scale projects in Nebraska?",//工程定制问题
                "acceptedAnswer": {
                    "@type": "Answer",
                    "text": "Yes, we specialize in custom engineering-grade cabinets for large-scale commercial and residential projects in Nebraska. Our team works closely with project managers to create custom designs, with lead times ranging from 2-4 weeks depending on project size."
                }
            }
        ]
    }
}
</script>
```

必选字段：批发覆盖区域、最小起订量、批发单价、对接联系方式；必选：工程批发FAQ（至少2组）、工程/服务权威资质（至少1项）；可选字段：批量折扣规则、区域供货周期、本地自提地点（批发自提专属）、付款条款说明。
核心说明：批发覆盖区域需明确（如“京津冀地区”“长三角核心城市”），对接联系方式需与区域匹配，便于本地商户联系；FAQ需聚焦工程客户高频问题（付款、定制、配送、安装），权威资质需与工程业务相关（如施工认证、商用产品认证）。

### 3.2 字段填写核心原则
1. 一致性原则：JSON-LD中的所有地理字段（地址、区域、经纬度、配送范围等）、FAQ字段、权威资质字段，必须与前端页面展示内容、后端接口返回数据完全一致，禁止出现数据冲突。
2. 完整性原则：必选字段不可缺失，FAQ模块核心页面（首页、商品详情页、批发页）至少包含2组问答，权威资质至少标注1项真实有效认证；可选字段按需补充，提升数据丰富度。
3. 规范性原则：地址格式统一为“省+市+区+详细地址”，国家代码统一填写“CN”/“US”；经纬度需精准，可通过地图工具获取；时间格式统一（如配送时效标注“周一至周日 09:00-18:00”）；FAQ问答需简洁明了，答案需包含地理/业务核心关键词；权威资质需标注名称、颁发机构、有效期（如有）。

### 3.3 禁止事项
- 禁止添加页面未展示的地理信息、FAQ内容、权威资质（如虚构发货仓、编造认证信息），否则可能被搜索引擎判定为作弊，降低网站权重。
- 禁止字段冗余或格式错误（如经纬度保留位数不足、地址缺失区域信息、FAQ问答格式不完整），避免搜索引擎解析失败。
- 禁止将JSON-LD放入异步加载的JavaScript文件中，搜索引擎可能无法抓取异步内容，导致结构化数据失效。
- 禁止FAQ问答堆砌关键词，需符合自然语言逻辑，与用户真实搜索意图匹配。

## 四、JSON-LD 及 GEO 信息适配页面清单（FAQ/权威信息适配）
| 页面类型                | 核心JSON-LD类型       | 适配内容                          |
|-------------------------|-----------------------|---------------------------------------|
| 首页                    | Organization+WebSite+CollectionPage | 企业权威资质、全站通用FAQ（地理/服务类） |
| 商品分类页（如橱柜页）  | CollectionPage        | 分类专属FAQ（商品+地理）、商品类目认证 |
| 零售商品详情页          | Product               | 商品专属FAQ（配送/自提/本地服务）、商品认证 |
| 批发商品详情页          | Product（嵌套Offer）  | 批发FAQ（付款/配送/工程定制）、工程资质 |
| 配送说明页              | DeliveryService       | 配送FAQ（时效/应急配送/安装）、配送服务资质 |
| 批发专区首页            | Product+Offer         | 批发通用FAQ、企业工程类权威资质        |
| 工程合作页              | Product+Service       | 工程对接FAQ、本地工程服务资质          |
| 联系页                  | Organization+ContactPoint | 本地联系方式FAQ、企业认证信息        |

## 五、GEO 优化补充措施（FAQ/权威信息优化）
### 5.1 页面内容优化
1. 标题与描述优化：严格沿用指定的H1和元描述，核心优化方向为在页面副标题、内容描述中补充地理关键词，如橱柜页副标题“xxx for 地理位置”，避免关键词堆砌；FAQ模块标题可融入地理关键词（如“xxx - FAQ for 地理位置&业务
2. 内容关联优化：在商品描述、服务说明中自然融入“地理位置+业务”等短语，贴合业务场景；FAQ模块需在页面可视化展示，与JSON-LD中的问答内容完全一致。
3. 图片优化：为产品图、案例图添加含地理关键词的alt标签，辅助搜索引擎识别区域属性；权威资质证书图片需添加alt标签，并在页面展示。
4. FAQ可视化优化：核心页面（首页、商品详情页、批发页）需在显著位置展示FAQ模块，问答内容与JSON-LD同步，样式适配移动端，提升用户体验；FAQ问题需包含地理/批发/商品等核心关键词，答案需准确且包含业务卖点。
5. 权威信息展示：企业/商品/服务的权威资质需在页面“关于我们”“商品详情”“服务说明”模块可视化展示，标注认证名称、颁发机构、有效期，与JSON-LD中的certification/hasCredential字段一致。

### 5.2 搜索引擎对接优化
1. Sitemap提交：按页面类型拆分Sitemap（如产品页Sitemap、功能页Sitemap），包含所有页面URL，标注lastmod字段（格式为YYYY-MM-DD），提交至谷歌搜索控制台、必应搜索资源平台（海外市场重点）/百度搜索资源平台（国内市场）；FAQ/权威信息的页面需优先标注lastmod，加速搜索引擎抓取。需包含的核心页面：首页、商品分类页面、商品详情页、FAQ聚合页（如有）。
2. 商户平台标注：在主流商户平台注册标注，包括谷歌商家档案（Google Business Profile）、Yelp、百度地图（针对华人客户），填写与网站一致的核心地址、电话及服务范围，标注企业权威资质，标注信息更新后同步更新网站Schema及页面内容。
3. 抓取权限优化：确保robots.txt未禁止抓取所有页面，动态URL（如分类页带参数URL）需实现静态化或添加canonical标签；核心页面（首页、橱柜分类页、联系页、FAQ页）添加内部链接，提升索引优先级，区域页面避免使用“noindex”标签。

### 动态URL（如分类页带参数URL）需实现静态化或添加canonical标签解释

例如
主分类页：https://xxx.com/category/framed
详细页：https://xxx.com/category/framed?pid=xxx    内容和主分类页高度相似，会被搜索引擎判定为重复内容，分散主分类页的权重。
解决
给所有动态参数 URL 添加 canonical 标签（JSON-LD外）
<head>
  <link rel="canonical" href="https://xxx.com/category/framed" />
</head>

## 六、验证与监控（FAQ/权威信息验证

### 6.1 JSON-LD有效性验证
1. 工具验证：页面上线后，使用谷歌Rich Results Test、百度搜索资源平台结构化数据验证工具，输入页面URL或JSON-LD内容，检查是否存在语法错误、字段缺失、数据冲突等问题；重点验证FAQPage、Certificate/OrganizationRole类型字段的解析状态，确保无报错。
2. 人工校验：逐一核对页面展示的地理信息、FAQ内容、权威资质与JSON-LD字段，确保一致性；检查JSON-LD是否成功渲染至页面头部，无异步加载问题；核对FAQ问答数量、权威资质数量是否满足规范要求。

### 6.2 效果监控与优化调整
1. 搜索控制台监控：通过百度搜索资源平台、谷歌搜索控制台，查看地理相关页面的索引情况、富媒体结果展示次数、点击量，分析JSON-LD是否生效；查看FAQ、权威信息相关的富摘要展示情况（如谷歌搜索结果中的FAQ折叠面板）。
2. 排名监控：定期跟踪核心地理关键词、权威/FAQ相关关键词的搜索排名，观察优化后排名变化趋势。
3. 数据调整：若排名无提升或JSON-LD解析失败，排查字段完整性、数据一致性、抓取权限等问题，针对性优化；若FAQ富摘要未展示，检查FAQPage字段格式、问答数量、内容相关性；根据业务变化（如认证、更新FAQ），及时更新JSON-LD及页面内容。

## 七、核心运维规范（含数据更新、Sitemap管理，FAQ/权威信息运维）
### 7.1 JSON-LD数据更新规范
当商品信息、地理数据、FAQ内容、权威资质发生变更时，需同步更新对应页面的JSON-LD标签，避免数据不一致导致搜索引擎信任度下降，具体更新场景及流程如下：

#### 7.1.1 核心更新场景
- 商品地理属性变更：服务区域扩展/缩减、区域安装服务调整、工程案例所在区域更新，需同步修正对应页面Schema中的“areaServed”字段。
- 页面信息变更：更新H1、元描述后，需同步调整Schema中的“name”“description”字段，确保一致性；商品SKU、名称变更需同步更新Product类型字段。
- 运营信息变更：核心地址、电话、营业时间变更，需同步更新Organization类型Schema及所有关联页面的展示内容、商户平台信息。
- FAQ内容变更：用户高频问题更新、业务规则调整（如配送时效变化），需同步更新FAQPage中的Question/Answer字段，确保与页面展示的FAQ一致。
- 权威资质变更：认证到期、认证/奖项、资质颁发机构信息变更，需同步更新certification/hasCredential字段，删除失效资质，有效资质。

#### 7.1.2 标准化更新流程
1. 后端同步更新：后端修改地理相关原始数据、FAQ问答库、权威资质库后，同步更新接口返回字段，确保字段名称、格式与原规范一致，字段需提前同步前端。
2. 前端标签更新：前端对接更新后的接口，重新渲染JSON-LD标签内容，重点核对一致性字段（如地址、经纬度、库存状态、FAQ问答、权威资质），避免拼写错误、格式偏差。
3. 验证环节：更新后通过谷歌Rich Results Test、百度结构化数据验证工具，逐项校验对应页面的JSON-LD有效性，重点验证FAQ、权威资质字段的解析状态，确认无解析错误。
4. 上线监控：页面上线后24-48小时内，在搜索控制台查看字段解析状态，确保更新后的数据正常被搜索引擎抓取；检查FAQ富摘要是否正常展示。


### 7.2 Sitemap更新管理
Sitemap需随地理相关页面增减、内容更新、FAQ/权威信息变更同步维护，确保搜索引擎及时抓取新页面、更新旧页面索引，具体规范如下：

#### 7.2.1 需更新Sitemap的场景
- 页面：添加博客页、FAQ聚合页等新页面时，需同步添加至Sitemap，补充博客页/FAQ页Schema（关联WebSite及Organization ID），融入地理关键词；页面需标注lastmod为当前日期。
- 页面变更：各页面H1、元描述、URL变更、FAQ、权威信息更新，需更新Sitemap中的URL及lastmod字段；页面删除需及时从Sitemap中移除，避免死链。
- 定期更新：建议每周全量更新一次Sitemap，重点核对所有页面的有效性、FAQ/权威信息字段的准确性，每月排查URL是否可达，移除失效页面。

#### 7.2.2 标准化更新流程
1. 生成Sitemap：后端或运维通过工具（如Go语言sitemap-go库）自动生成包含所有地理相关页面、FAQ页、权威信息页的Sitemap，格式为XML，确保每个URL标注最后更新时间（lastmod字段）。
2. 提交更新：将生成的Sitemap提交至百度搜索资源平台、谷歌搜索控制台，百度支持手动提交或通过API批量提交，谷歌可关联站点自动更新；FAQ/权威信息更新的页面需单独标记，优先提交。
3. 索引监控：提交后1-3天内，在搜索控制台查看Sitemap状态，确认无错误（如URL不可达、格式错误），跟踪页面索引进度；查看FAQ/权威信息相关页面的索引状态。

#### 7.2.3 注意事项
Sitemap中仅包含可索引页面，禁止添加被robots禁止抓取、带noindex标签的页面；地理相关页面、FAQ页、权威信息页需单独分类标注（可按页面类型拆分Sitemap，如商品页、批发页、配送页、FAQ页），提升抓取优先级。