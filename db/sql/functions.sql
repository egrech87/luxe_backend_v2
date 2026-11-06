CREATE OR REPLACE FUNCTION is_aperture_combined(apertureId bigint)
RETURNS boolean AS $isCombined$
declare
	isCombined boolean;
	total integer;
BEGIN
   SELECT count(*) into total FROM apertures WHERE apertures.parent_aperture_id = apertureId AND apertures.deleted_at IS NULL;
   IF total > 0 THEN
		isCombined = true;
   ELSE
		isCombined = false;
   END IF;
   
   RETURN isCombined;
END;
$isCombined$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION total_apertures_in_quote(quoteId bigint)
RETURNS integer AS $total$
declare
	total integer;
BEGIN
   SELECT COALESCE(SUM(apertures.quantity), 0) into total
   FROM quotes INNER JOIN apertures ON quotes.id = apertures.quote_id 
   WHERE quotes.id = quoteId AND apertures.parent_aperture_id IS null AND apertures.deleted_at IS NULL;
   
   RETURN total;
END;
$total$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_friendly_name(apertureId BIGINT)
RETURNS TEXT AS $friendlyName$
declare
	apertureFriendlyName TEXT;
	profileFriendlyName TEXT;
	isFriendlyNameCustom BOOLEAN;
BEGIN
   	SELECT a.friendly_name, a.is_friendly_name_custom, p.friendly_name 
	INTO apertureFriendlyName, isFriendlyNameCustom, profileFriendlyName
	FROM apertures a LEFT OUTER JOIN profiles p ON aperture_profile_id(apertureId) = p.id
	WHERE a.id=apertureId;
	
	IF isFriendlyNameCustom THEN
		RETURN apertureFriendlyName;
	END IF;
   
   	RETURN profileFriendlyName;
END;
$friendlyName$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_profile_id(apertureId bigint)
RETURNS numeric AS $profileId$
declare
	parentApertureId bigint;
	profileId numeric;
	isProfileIdCustom BOOLEAN;
BEGIN
   	SELECT apertures.profile_id, apertures.is_profile_id_custom, apertures.parent_aperture_id 
	INTO profileId, isProfileIdCustom, parentApertureId
	FROM apertures
	WHERE apertures.id=apertureId;
	
	IF isProfileIdCustom OR parentApertureId IS NULL THEN
		RETURN profileId;
	END IF;
	
	-- Get the Profile ID from the parent aperture
   	SELECT parent.profile_id 
	INTO profileId
	FROM apertures AS parent INNER JOIN apertures AS child ON parent.id=child.parent_aperture_id
	WHERE child.id=apertureId;
   
   	RETURN profileId;
END;
$profileId$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_color_id(apertureId bigint)
RETURNS numeric AS $colorId$
declare
	parentApertureId bigint;
	colorId numeric;
	isColorIdCustom BOOLEAN;
BEGIN
   	SELECT apertures.color_id, apertures.is_color_id_custom, apertures.parent_aperture_id INTO colorId, isColorIdCustom, parentApertureId
	FROM apertures
	WHERE apertures.id=apertureId;
	
	IF isColorIdCustom OR parentApertureId IS NULL THEN
		RETURN colorId;
	END IF;
	
	-- Get the Color ID from the parent aperture
   	SELECT parent.color_id INTO colorId
	FROM apertures AS parent INNER JOIN apertures AS child ON parent.id=child.parent_aperture_id
	WHERE child.id=apertureId;
   
   	RETURN colorId;
END;
$colorId$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_glass_id(apertureId bigint)
RETURNS numeric AS $glassId$
declare
	parentApertureId bigint;
	glassId numeric;
	isGlassIdCustom BOOLEAN;
BEGIN
   	SELECT apertures.glass_id, apertures.is_glass_id_custom, apertures.parent_aperture_id INTO glassId, isGlassIdCustom, parentApertureId
	FROM apertures
	WHERE apertures.id=apertureId;
	
	IF isGlassIdCustom OR parentApertureId IS NULL THEN
		RETURN glassId;
	END IF;
	
	-- Get the Glass ID from the parent aperture
   	SELECT parent.glass_id INTO glassId
	FROM apertures AS parent INNER JOIN apertures AS child ON parent.id=child.parent_aperture_id
	WHERE child.id=apertureId;
   
   	RETURN glassId;
END;
$glassId$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_type(apertureId bigint)
RETURNS numeric AS $type$
declare
	parentApertureId bigint;
	_type numeric;
	isTypeCustom BOOLEAN;
BEGIN
   	SELECT apertures.type, apertures.is_type_custom, apertures.parent_aperture_id 
	INTO _type, isTypeCustom, parentApertureId
	FROM apertures
	WHERE apertures.id=apertureId;
	
	IF isTypeCustom OR parentApertureId IS NULL THEN
		RETURN _type;
	END IF;
	
	-- Get the Type from the parent aperture
   	SELECT parent.type 
	INTO _type
	FROM apertures AS parent INNER JOIN apertures AS child ON parent.id=child.parent_aperture_id
	WHERE child.id=apertureId;
   
   	RETURN _type;
END;
$type$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_handle(apertureId bigint)
RETURNS numeric AS $handle$
declare
	parentApertureId bigint;
	handle numeric;
	isHandleCustom BOOLEAN;
BEGIN
   	SELECT apertures.handle, apertures.is_handle_custom, apertures.parent_aperture_id 
	INTO handle, isHandleCustom, parentApertureId
	FROM apertures
	WHERE apertures.id=apertureId;
	
	IF isHandleCustom OR parentApertureId IS NULL THEN
		RETURN handle;
	END IF;
	
	-- Get the Handle from the parent aperture
   	SELECT parent.handle 
	INTO handle
	FROM apertures AS parent INNER JOIN apertures AS child ON parent.id=child.parent_aperture_id
	WHERE child.id=apertureId;
   
   	RETURN handle;
END;
$handle$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_area(apertureId bigint)
RETURNS numeric AS $area$
declare
	area numeric;
	customArea numeric;
	isAreaCustom BOOLEAN;
BEGIN
   	SELECT COALESCE(apertures.width, 0) * COALESCE(apertures.height, 0) / 1000000.0, apertures.area, apertures.is_area_custom into area, customArea, isAreaCustom
	FROM apertures 
	WHERE apertures.id = apertureId;

	IF isAreaCustom THEN
		RETURN COALESCE(customArea, 0);
	END IF;
   
   RETURN area;
END;
$area$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_color_price_factor(apertureId bigint)
RETURNS numeric AS $colorPriceFactor$
declare
	parentApertureId bigint;
	colorPriceFactor numeric;
	customPriceFactor numeric;
	isPriceFactorCustom BOOLEAN;
	isColorIdCustom BOOLEAN;
BEGIN
   	SELECT colors.price_factor, apertures.color_price_factor, apertures.is_color_price_factor_custom, apertures.is_color_id_custom, apertures.parent_aperture_id INTO colorPriceFactor, customPriceFactor, isPriceFactorCustom, isColorIdCustom, parentApertureId
	FROM colors INNER JOIN apertures ON colors.id = aperture_color_id(apertureId)
	WHERE apertures.id=apertureId;
	
	IF isPriceFactorCustom THEN
		RETURN COALESCE(customPriceFactor, 0);
	END IF;
	
	IF isColorIdCustom OR parentApertureId IS NULL THEN
		RETURN COALESCE(colorPriceFactor, 0);
	END IF;
	
	-- Get the Price factor from the parent aperture
   	SELECT colors.price_factor, parent.color_price_factor, parent.is_color_price_factor_custom INTO colorPriceFactor, customPriceFactor, isPriceFactorCustom
	FROM colors INNER JOIN apertures AS parent ON colors.id = aperture_color_id(parent.id) INNER JOIN apertures AS child ON parent.id=child.parent_aperture_id
	WHERE child.id=apertureId;
	
	IF isPriceFactorCustom THEN
		RETURN COALESCE(customPriceFactor, 0);
	END IF;
   
   	RETURN COALESCE(colorPriceFactor, 0);
END;
$colorPriceFactor$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_glass_price_factor(apertureId bigint)
RETURNS numeric AS $glassPriceFactor$
declare
	parentApertureId bigint;
	glassPriceFactor numeric;
	customPriceFactor numeric;
	isPriceFactorCustom BOOLEAN;
	isGlassIdCustom BOOLEAN;
BEGIN
   	SELECT glasses.price_factor, apertures.glass_price_factor, apertures.is_glass_price_factor_custom, apertures.is_glass_id_custom, apertures.parent_aperture_id INTO glassPriceFactor, customPriceFactor, isPriceFactorCustom, isGlassIdCustom, parentApertureId
	FROM glasses INNER JOIN apertures ON glasses.id = aperture_glass_id(apertureId)
	WHERE apertures.id=apertureId;
	
	IF isPriceFactorCustom THEN
		RETURN COALESCE(customPriceFactor, 0);
	END IF;
	
	IF isGlassIdCustom OR parentApertureId IS NULL THEN
		RETURN COALESCE(glassPriceFactor, 0);
	END IF;
	
	-- Get the Price factor from the parent aperture
   	SELECT glasses.price_factor, parent.glass_price_factor, parent.is_glass_price_factor_custom INTO glassPriceFactor, customPriceFactor, isPriceFactorCustom
	FROM glasses INNER JOIN apertures AS parent ON glasses.id = aperture_glass_id(parent.id) INNER JOIN apertures AS child ON parent.id=child.parent_aperture_id
	WHERE child.id=apertureId;
	
	IF isPriceFactorCustom THEN
		RETURN COALESCE(customPriceFactor, 0);
	END IF;
   
   	RETURN COALESCE(glassPriceFactor, 0);
END;
$glassPriceFactor$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_matrix_cost(apertureId bigint)
RETURNS numeric AS $matrixCost$
declare
	profileId bigint;
	width TEXT;
	height TEXT;
	matrixCost numeric;
	sumMatrixCost numeric;
	customMatrixCost numeric;
	isMatrixCostCustom BOOLEAN;
	isCombined BOOLEAN;
	priceMatrix TEXT;
   	_key   text;
   	_value text;
	chosenWidth TEXT;
	chosenHeight TEXT;
BEGIN
   	SELECT CAST(COALESCE(apertures.width, 0) AS TEXT), CAST(COALESCE(apertures.height, 0) AS TEXT), COALESCE(apertures.matrix_cost, 0), apertures.is_matrix_cost_custom, aperture_profile_id(apertureId), is_aperture_combined(apertureId)
	INTO width, height, customMatrixCost, isMatrixCostCustom, profileId, isCombined
	FROM apertures 
	WHERE apertures.id = apertureId;
	
	IF isCombined THEN
		SELECT SUM(aperture_matrix_cost(id) * quantity) into sumMatrixCost
		FROM apertures 
		WHERE apertures.parent_aperture_id = apertureId AND apertures.deleted_at IS NULL;

		RETURN sumMatrixCost;
	END IF;

	IF isMatrixCostCustom THEN
		RETURN COALESCE(customMatrixCost, 0);
	END IF;
	
	IF profileId IS NULL OR width = '0' OR height = '0' THEN
		RETURN 0;
	END IF;

	SELECT profiles.price_matrix into priceMatrix 
	FROM profiles 
	WHERE profiles.id = profileId;

	-- Choose the width
	FOR _key, _value IN
       	SELECT * FROM json_each(priceMatrix::json)
    LOOP
       	chosenWidth = _key;
	   
	   	IF chosenWidth::int >= width::int THEN
	   		EXIT;
		END IF;
    END LOOP;
	
	-- Choose the height
	FOR _key, _value IN
       	SELECT * FROM json_each(priceMatrix::json->chosenWidth)
    LOOP
       	chosenHeight = _key;
	   
	   	IF chosenHeight::int >= height::int THEN
	   		EXIT;
		END IF;
    END LOOP;
	
	-- Choose the matrix cost
 	matrixCost = priceMatrix::json->chosenWidth->chosenHeight;

	RETURN COALESCE(matrixCost, 0);
END;
$matrixCost$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_color_cost(apertureId bigint)
RETURNS numeric AS $colorCost$
declare
	sumColorCost numeric;
	customColorCost numeric;
	matrixCost numeric;
	colorPriceFactor numeric;
	isColorCostCustom BOOLEAN;
	isCombined BOOLEAN;
BEGIN
   	SELECT COALESCE(apertures.color_cost, 0), aperture_matrix_cost(apertureid), aperture_color_price_factor(apertureid), apertures.is_color_cost_custom, is_aperture_combined(apertureId) 
	into customColorCost, matrixCost, colorPriceFactor, isColorCostCustom, isCombined
	FROM apertures 
	WHERE apertures.id = apertureId;
	
	IF isCombined THEN
		SELECT SUM(aperture_color_cost(id) * quantity) into sumColorCost
		FROM apertures 
		WHERE apertures.parent_aperture_id = apertureId AND apertures.deleted_at IS NULL;

		RETURN sumColorCost;
	END IF;

	IF isColorCostCustom THEN
		RETURN COALESCE(customColorCost, 0);
	END IF;

	RETURN matrixCost * (colorPriceFactor - 1);
END;
$colorCost$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_glass_cost(apertureId bigint)
RETURNS numeric AS $glassCost$
declare
	sumGlassCost numeric;
	customGlassCost numeric;
	matrixCost numeric;
	glassPriceFactor numeric;
	isGlassCostCustom BOOLEAN;
	isCombined BOOLEAN;
BEGIN
   	SELECT COALESCE(apertures.glass_cost, 0), aperture_matrix_cost(apertureid), aperture_glass_price_factor(apertureid), apertures.is_glass_cost_custom, is_aperture_combined(apertureId) 
	into customGlassCost, matrixCost, glassPriceFactor, isGlassCostCustom, isCombined
	FROM apertures 
	WHERE apertures.id = apertureId;
	
	IF isCombined THEN
		SELECT SUM(aperture_glass_cost(id) * quantity) into sumGlassCost
		FROM apertures 
		WHERE apertures.parent_aperture_id = apertureId AND apertures.deleted_at IS NULL;

		RETURN sumGlassCost;
	END IF;

	IF isGlassCostCustom THEN
		RETURN COALESCE(customGlassCost, 0);
	END IF;

	RETURN matrixCost * (glassPriceFactor - 1);
END;
$glassCost$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_nets_quantity(apertureId bigint)
RETURNS numeric AS $netsQuantity$
declare
	netsQuantity numeric;
	sumNetsQuantity numeric;
	isCombined BOOLEAN;
BEGIN
   	SELECT COALESCE(apertures.nets_quantity, 0), is_aperture_combined(apertureId) into netsQuantity, isCombined
	FROM apertures 
	WHERE apertures.id = apertureId;

	IF isCombined THEN
		SELECT SUM(aperture_nets_quantity(id) * quantity) into sumNetsQuantity
		FROM apertures 
		WHERE apertures.parent_aperture_id = apertureId AND apertures.deleted_at IS NULL;

		RETURN sumNetsQuantity;
	END IF;
	
	RETURN netsQuantity;
END;
$netsQuantity$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_nets_price(apertureId bigint)
RETURNS numeric AS $netsPrice$
declare
	parentApertureId bigint;
	netsPrice numeric;
	isNetsPriceCustom BOOLEAN;
BEGIN
   	SELECT apertures.nets_price, apertures.is_nets_price_custom, apertures.parent_aperture_id INTO netsPrice, isNetsPriceCustom, parentApertureId
	FROM apertures
	WHERE apertures.id=apertureId;
	
	IF isNetsPriceCustom OR parentApertureId IS NULL THEN
		RETURN netsPrice;
	END IF;
	
	-- Get the Net Price from the parent aperture
   	SELECT parent.nets_price INTO netsPrice
	FROM apertures AS parent INNER JOIN apertures AS child ON parent.id=child.parent_aperture_id
	WHERE child.id=apertureId;
   
   	RETURN netsPrice;
END;
$netsPrice$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_nets_cost(apertureId bigint)
RETURNS numeric AS $netsCost$
declare
	sumNetsCost numeric;
	customNetsCost numeric;
	netsQuantity numeric;
	netsPrice numeric;
	isNetsCostCustom BOOLEAN;
	isCombined BOOLEAN;
BEGIN
   	SELECT COALESCE(apertures.nets_cost, 0), aperture_nets_quantity(apertureid), aperture_nets_price(apertureid), apertures.is_nets_cost_custom, is_aperture_combined(apertureId) 
	INTO customNetsCost, netsQuantity, netsPrice, isNetsCostCustom, isCombined
	FROM apertures 
	WHERE apertures.id = apertureId;
	
	IF isCombined THEN
		SELECT SUM(aperture_nets_cost(id) * quantity) into sumNetsCost
		FROM apertures 
		WHERE apertures.parent_aperture_id = apertureId AND apertures.deleted_at IS NULL;

		RETURN sumNetsCost;
	END IF;

	IF isNetsCostCustom THEN
		RETURN COALESCE(customNetsCost, 0);
	END IF;

	RETURN netsQuantity * netsPrice;
END;
$netsCost$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_other_cost(apertureId bigint)
RETURNS numeric AS $otherCost$
declare
	otherCost numeric;
	sumOtherCost numeric;
	isCombined BOOLEAN;
BEGIN
   	SELECT COALESCE(apertures.other_cost, 0), is_aperture_combined(apertureId) 
	INTO otherCost, isCombined
	FROM apertures 
	WHERE apertures.id = apertureId;

	IF isCombined THEN
		SELECT SUM(aperture_other_cost(id) * quantity) into sumOtherCost
		FROM apertures 
		WHERE apertures.parent_aperture_id = apertureId AND apertures.deleted_at IS NULL;

		RETURN sumOtherCost;
	END IF;
	
	RETURN otherCost;
END;
$otherCost$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_other_cost_comment(apertureId bigint)
RETURNS text AS $otherCostComment$
declare
	combinedOtherCostComment text;
	otherCostComment text;
	isCombined BOOLEAN;
BEGIN
   	SELECT apertures.other_cost_comment, is_aperture_combined(apertureId)
	into otherCostComment, isCombined
	FROM apertures 
	WHERE apertures.id = apertureId;
	
	IF isCombined THEN
		SELECT string_agg(tmp.other_cost_comment, E'\n')
		INTO combinedOtherCostComment
		FROM (
			SELECT apertures.other_cost_comment
			FROM apertures
			WHERE apertures.parent_aperture_id = apertureId AND apertures.other_cost_comment != '' AND apertures.deleted_at IS NULL
			ORDER BY apertures.id ASC
		) AS tmp;

		RETURN COALESCE(combinedOtherCostComment, '');
	END IF;

	RETURN otherCostComment;
END;
$otherCostComment$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_supplier_cost(apertureId bigint)
RETURNS numeric AS $supplierCost$
declare
	sumSupplierCost numeric;
	customSupplierCost numeric;
	colorCost numeric;
	glassCost numeric;
	matrixCost numeric;
	netsCost numeric;
	otherCost numeric;
	isSupplierCostCustom BOOLEAN;
	isCombined BOOLEAN;
BEGIN
   	SELECT COALESCE(apertures.supplier_cost, 0), aperture_color_cost(apertureid), aperture_glass_cost(apertureid), aperture_matrix_cost(apertureid), aperture_nets_cost(apertureid), aperture_other_cost(apertureid), apertures.is_supplier_cost_custom, is_aperture_combined(apertureId) 
	INTO customSupplierCost, colorCost, glassCost, matrixCost, netsCost, otherCost, isSupplierCostCustom, isCombined
	FROM apertures 
	WHERE apertures.id = apertureId;
	
	IF isCombined THEN
		SELECT SUM(aperture_supplier_cost(id) * quantity) 
		INTO sumSupplierCost
		FROM apertures 
		WHERE apertures.parent_aperture_id = apertureId AND apertures.deleted_at IS NULL;

		RETURN sumSupplierCost;
	END IF;

	IF isSupplierCostCustom THEN
		RETURN COALESCE(customSupplierCost, 0);
	END IF;

	RETURN colorCost + glassCost + matrixCost + netsCost + otherCost;
END;
$supplierCost$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_shipping_cost(apertureId bigint)
RETURNS numeric AS $shippingCost$
declare
	sumShippingCost numeric;
	customShippingCost numeric;
	shippingPerUnit numeric;
	area numeric;
	isShippingCostCustom BOOLEAN;
	isCombined BOOLEAN;
BEGIN
   	SELECT COALESCE(apertures.shipping_cost, 0), COALESCE(quotes.shipping_per_unit, 0), aperture_area(apertureid), apertures.is_shipping_cost_custom, is_aperture_combined(apertureId) 
	INTO customShippingCost, shippingPerUnit, area, isShippingCostCustom, isCombined
	FROM apertures INNER JOIN quotes ON quotes.id = apertures.quote_id
	WHERE apertures.id = apertureId;
	
	IF isCombined THEN
		SELECT SUM(aperture_shipping_cost(id) * quantity) 
		INTO sumShippingCost
		FROM apertures 
		WHERE apertures.parent_aperture_id = apertureId AND apertures.deleted_at IS NULL;

		RETURN sumShippingCost;
	END IF;

	IF isShippingCostCustom THEN
		RETURN COALESCE(customShippingCost, 0);
	END IF;

	RETURN shippingPerUnit * area;
END;
$shippingCost$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_installation_per_unit(apertureId bigint)
RETURNS numeric AS $installationPerUnit$
declare
	parentApertureId bigint;
	quoteInstallationPerUnit numeric;
	customInstallationPerUnit numeric;
	isInstallationPerUnitCustom BOOLEAN;
BEGIN
   	SELECT quotes.installation_per_unit, apertures.installation_per_unit, apertures.is_installation_per_unit_custom, apertures.parent_aperture_id 
	INTO quoteInstallationPerUnit, customInstallationPerUnit, isInstallationPerUnitCustom, parentApertureId
	FROM quotes INNER JOIN apertures ON quotes.id = apertures.quote_id
	WHERE apertures.id=apertureId;
	
	IF isInstallationPerUnitCustom THEN
		RETURN COALESCE(customInstallationPerUnit, 0);
	END IF;
	
	IF parentApertureId IS NULL THEN
		RETURN COALESCE(quoteInstallationPerUnit, 0);
	END IF;
	
	-- Get the Installation per Unit from the parent aperture
   	SELECT parent.installation_per_unit, parent.is_installation_per_unit_custom 
	INTO customInstallationPerUnit, isInstallationPerUnitCustom
	FROM apertures AS parent INNER JOIN apertures AS child ON parent.id=child.parent_aperture_id
	WHERE child.id=apertureId;
	
	IF isInstallationPerUnitCustom THEN
		RETURN COALESCE(customInstallationPerUnit, 0);
	END IF;
   
   	RETURN COALESCE(quoteInstallationPerUnit, 0);
END;
$installationPerUnit$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_installation_cost(apertureId bigint)
RETURNS numeric AS $installationCost$
declare
	sumInstallationCost numeric;
	customInstallationCost numeric;
	installationPerUnit numeric;
	area numeric;
	isInstallationCostCustom BOOLEAN;
	isCombined BOOLEAN;
BEGIN
   	SELECT COALESCE(apertures.installation_cost, 0), aperture_installation_per_unit(apertureid), aperture_area(apertureid), apertures.is_installation_cost_custom, is_aperture_combined(apertureId) 
	INTO customInstallationCost, installationPerUnit, area, isInstallationCostCustom, isCombined
	FROM apertures
	WHERE apertures.id = apertureId;
	
	IF isCombined THEN
		SELECT SUM(aperture_installation_cost(id) * quantity) 
		INTO sumInstallationCost
		FROM apertures 
		WHERE apertures.parent_aperture_id = apertureId AND apertures.deleted_at IS NULL;

		RETURN sumInstallationCost;
	END IF;

	IF isInstallationCostCustom THEN
		RETURN COALESCE(customInstallationCost, 0);
	END IF;

	RETURN installationPerUnit * area;
END;
$installationCost$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_contingency_cost(apertureId bigint)
RETURNS numeric AS $contingencyCost$
declare
	sumContingencyCost numeric;
	customContingencyCost numeric;
	contingencyPerc numeric;
	shippingCost numeric;
	installationCost numeric;
	supplierCost numeric;
	isContingencyCostCustom BOOLEAN;
	isCombined BOOLEAN;
BEGIN
   	SELECT COALESCE(apertures.contingency_cost, 0), COALESCE(quotes.contingency_perc, 0), aperture_shipping_cost(apertureid), aperture_installation_cost(apertureid), aperture_supplier_cost(apertureid), apertures.is_contingency_cost_custom, is_aperture_combined(apertureId) 
	INTO customContingencyCost, contingencyPerc, shippingCost, installationCost, supplierCost, isContingencyCostCustom, isCombined
	FROM apertures INNER JOIN quotes ON quotes.id = apertures.quote_id
	WHERE apertures.id = apertureId;
	
	IF isCombined THEN
		SELECT SUM(aperture_contingency_cost(id) * quantity) 
		INTO sumContingencyCost
		FROM apertures 
		WHERE apertures.parent_aperture_id = apertureId AND apertures.deleted_at IS NULL;

		RETURN sumContingencyCost;
	END IF;

	IF isContingencyCostCustom THEN
		RETURN COALESCE(customContingencyCost, 0);
	END IF;

	RETURN (shippingCost + installationCost + supplierCost) * contingencyPerc / 100.0;
END;
$contingencyCost$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_item_cost(apertureId bigint)
RETURNS numeric AS $itemCost$
declare
	sumItemCost numeric;
	customItemCost numeric;
	shippingCost numeric;
	installationCost numeric;
	supplierCost numeric;
	contingencyCost numeric;
	isItemCostCustom BOOLEAN;
	isCombined BOOLEAN;
BEGIN
   	SELECT COALESCE(apertures.item_cost, 0), aperture_shipping_cost(apertureid), aperture_installation_cost(apertureid), aperture_supplier_cost(apertureid), aperture_contingency_cost(apertureid), apertures.is_item_cost_custom, is_aperture_combined(apertureId) 
	INTO customItemCost, shippingCost, installationCost, supplierCost, contingencyCost, isItemCostCustom, isCombined
	FROM apertures
	WHERE apertures.id = apertureId;
	
	IF isCombined THEN
		SELECT SUM(aperture_total_cost(id)) 
		INTO sumItemCost
		FROM apertures 
		WHERE apertures.parent_aperture_id = apertureId AND apertures.deleted_at IS NULL;

		RETURN sumItemCost;
	END IF;

	IF isItemCostCustom THEN
		RETURN COALESCE(customItemCost, 0);
	END IF;

	RETURN shippingCost + installationCost + supplierCost + contingencyCost;
END;
$itemCost$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_profit_margin(apertureId bigint)
RETURNS numeric AS $profitMargin$
declare
	parentApertureId bigint;
	quoteProfitMargin numeric;
	customProfitMargin numeric;
	isProfitMarginCustom BOOLEAN;
BEGIN
   	SELECT quotes.profit_margin, apertures.profit_margin, apertures.is_profit_margin_custom, apertures.parent_aperture_id 
	INTO quoteProfitMargin, customProfitMargin, isProfitMarginCustom, parentApertureId
	FROM quotes INNER JOIN apertures ON quotes.id = apertures.quote_id
	WHERE apertures.id=apertureId;
	
	IF isProfitMarginCustom THEN
		RETURN COALESCE(customProfitMargin, 0);
	END IF;
	
	IF parentApertureId IS NULL THEN
		RETURN COALESCE(quoteProfitMargin, 0);
	END IF;
	
	-- Get the Profit Margin from the parent aperture
   	SELECT parent.profit_margin, parent.is_profit_margin_custom 
	INTO customProfitMargin, isProfitMarginCustom
	FROM apertures AS parent INNER JOIN apertures AS child ON parent.id=child.parent_aperture_id
	WHERE child.id=apertureId;
	
	IF isProfitMarginCustom THEN
		RETURN COALESCE(customProfitMargin, 0);
	END IF;
   
   	RETURN COALESCE(quoteProfitMargin, 0);
END;
$profitMargin$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_selling_price(apertureId bigint)
RETURNS numeric AS $sellingPrice$
declare
	sumSellingPrice numeric;
	customSellingPrice numeric;
	itemCost numeric;
	profitMargin numeric;
	isSellingPriceCustom BOOLEAN;
	isCombined BOOLEAN;
BEGIN
   	SELECT COALESCE(apertures.selling_price, 0), aperture_item_cost(apertureid), aperture_profit_margin(apertureid), apertures.is_selling_price_custom, is_aperture_combined(apertureId) 
	INTO customSellingPrice, itemCost, profitMargin, isSellingPriceCustom, isCombined
	FROM apertures
	WHERE apertures.id = apertureId;
	
	IF isCombined THEN
		SELECT SUM(aperture_total_selling_price(id)) 
		INTO sumSellingPrice
		FROM apertures 
		WHERE apertures.parent_aperture_id = apertureId AND apertures.deleted_at IS NULL;

		RETURN sumSellingPrice;
	END IF;

	IF isSellingPriceCustom THEN
		RETURN COALESCE(customSellingPrice, 0);
	END IF;

	RETURN itemCost / (100.0 - profitMargin) * 100;
END;
$sellingPrice$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_discounted_selling_price(apertureId bigint)
RETURNS numeric AS $discountedSellingPrice$
declare
	sumDiscountedSellingPrice numeric;
	customDiscountedSellingPrice numeric;
	sellingPrice numeric;
	discountPercentage numeric;
	isDiscountedSellingPriceCustom BOOLEAN;
	isCombined BOOLEAN;
BEGIN
   	SELECT COALESCE(apertures.discounted_selling_price, 0), aperture_selling_price(apertureid), quotes.discount_perc, apertures.is_discounted_selling_price_custom, is_aperture_combined(apertureId) 
	INTO customDiscountedSellingPrice, sellingPrice, discountPercentage, isDiscountedSellingPriceCustom, isCombined
	FROM apertures INNER JOIN quotes ON quotes.id = apertures.quote_id
	WHERE apertures.id = apertureId;
	
	IF isCombined THEN
		SELECT SUM(aperture_total_discounted_selling_price(id)) 
		INTO sumDiscountedSellingPrice
		FROM apertures 
		WHERE apertures.parent_aperture_id = apertureId AND apertures.deleted_at IS NULL;

		RETURN sumDiscountedSellingPrice;
	END IF;

	IF isDiscountedSellingPriceCustom THEN
		RETURN COALESCE(customDiscountedSellingPrice, 0);
	END IF;

	RETURN sellingPrice * (100 - discountPercentage) / 100.0;
END;
$discountedSellingPrice$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_image_id(apertureId bigint)
RETURNS integer AS $imageId$
declare
	profileId bigint;
	imageId bigint;
	leftImageId bigint;
	rightImageId bigint;
	handle integer;
	isImageIdCustom boolean;
BEGIN
	SELECT aperture_profile_id(apertureId), aperture_handle(apertureId), apertures.image_id, apertures.is_image_id_custom
	INTO profileId, handle, imageId, isImageIdCustom
   	FROM apertures 
   	WHERE apertures.id = apertureId;
   
   	IF isImageIdCustom THEN
		RETURN imageId;
   	END IF;
   
   	IF profileId IS NULL OR handle IS NULL THEN
   		RETURN NULL;
	END IF;
   
   	SELECT profiles.left_image_id, profiles.right_image_id
   	INTO leftImageId, rightImageId
   	FROM profiles 
 	WHERE profiles.id = profileId;
   
   	IF handle = 0 THEN
		RETURN leftImageId;
	ELSEIF handle = 1 THEN
		RETURN rightImageId;
	ELSE 
		RETURN NULL;
	END IF;
END;
$imageId$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_total_area(apertureId bigint)
RETURNS numeric AS $totalArea$
declare
	totalArea numeric;
	customTotalArea numeric;
	isTotalAreaCustom BOOLEAN;
BEGIN
   	SELECT aperture_area(apertureId) * COALESCE(apertures.quantity, 0), apertures.total_area, apertures.is_total_area_custom 
	INTO totalArea, customTotalArea, isTotalAreaCustom
	FROM apertures 
	WHERE apertures.id = apertureId;

	IF isTotalAreaCustom THEN
		RETURN COALESCE(customTotalArea, 0);
	END IF;
   
   RETURN totalArea;
END;
$totalArea$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_total_selling_price(apertureId bigint)
RETURNS numeric AS $totalSellingPrice$
declare
	totalSellingPrice numeric;
	customTotalSellingPrice numeric;
	isTotalSellingPriceCustom BOOLEAN;
BEGIN
   	SELECT aperture_selling_price(apertureId) * COALESCE(apertures.quantity, 0), apertures.total_selling_price, apertures.is_total_selling_price_custom 
	INTO totalSellingPrice, customTotalSellingPrice, isTotalSellingPriceCustom
	FROM apertures 
	WHERE apertures.id = apertureId;

	IF isTotalSellingPriceCustom THEN
		RETURN COALESCE(customTotalSellingPrice, 0);
	END IF;
   
   RETURN totalSellingPrice;
END;
$totalSellingPrice$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_total_discounted_selling_price(apertureId bigint)
RETURNS numeric AS $totalDiscountedSellingPrice$
declare
	totalDiscountedSellingPrice numeric;
	customTotalDiscountedSellingPrice numeric;
	isTotalDiscountedSellingPriceCustom BOOLEAN;
BEGIN
   	SELECT aperture_discounted_selling_price(apertureId) * COALESCE(apertures.quantity, 0), apertures.total_discounted_selling_price, apertures.is_total_discounted_selling_price_custom 
	INTO totalDiscountedSellingPrice, customTotalDiscountedSellingPrice, isTotalDiscountedSellingPriceCustom
	FROM apertures 
	WHERE apertures.id = apertureId;

	IF isTotalDiscountedSellingPriceCustom THEN
		RETURN COALESCE(customTotalDiscountedSellingPrice, 0);
	END IF;
   
   RETURN totalDiscountedSellingPrice;
END;
$totalDiscountedSellingPrice$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_total_cost(apertureId bigint)
RETURNS numeric AS $totalCost$
declare
	totalCost numeric;
	customTotalCost numeric;
	isTotalCostCustom BOOLEAN;
BEGIN
   	SELECT aperture_item_cost(apertureId) * COALESCE(apertures.quantity, 0), apertures.total_cost, apertures.is_total_cost_custom 
	INTO totalCost, customTotalCost, isTotalCostCustom
	FROM apertures 
	WHERE apertures.id = apertureId;

	IF isTotalCostCustom THEN
		RETURN COALESCE(customTotalCost, 0);
	END IF;
   
   RETURN totalCost;
END;
$totalCost$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_total_profit(apertureId bigint)
RETURNS numeric AS $totalProfit$
declare
	totalDiscountedSellingPrice numeric;
	totalCost numeric;
BEGIN
   	SELECT aperture_total_discounted_selling_price(apertureid), aperture_total_cost(apertureid)
	INTO totalDiscountedSellingPrice, totalCost
	FROM apertures
	WHERE apertures.id = apertureId;

	RETURN totalDiscountedSellingPrice - totalCost;
END;
$totalProfit$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_updated_at(apertureId bigint)
RETURNS TIMESTAMP WITH TIME ZONE AS $updatedAt$
declare
	updatedAt TIMESTAMP WITH TIME ZONE;
BEGIN
   	SELECT MAX(apertures.updated_at)
	INTO updatedAt
	FROM apertures
	WHERE apertures.id = apertureId OR (apertures.parent_aperture_id = apertureId AND apertures.deleted_at IS NULL);

	RETURN updatedAt;
END;
$updatedAt$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_updated_by(apertureId bigint)
RETURNS BIGINT AS $updatedBy$
declare
	updatedBy BIGINT;
BEGIN
   	SELECT apertures.updated_by
	INTO updatedBy
	FROM apertures
	WHERE apertures.updated_at = aperture_updated_at(apertureId) AND (apertures.id = apertureId OR (apertures.parent_aperture_id = apertureId AND apertures.deleted_at IS NULL))
	LIMIT 1;

	RETURN updatedBy;
END;
$updatedBy$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_glass_area(apertureId bigint)
RETURNS NUMERIC AS $glassArea$
declare
	width NUMERIC;
	height NUMERIC;
	glassArea NUMERIC;
BEGIN
   SELECT apertures.width, apertures.height
   INTO width, height 
   FROM apertures WHERE apertures.id = apertureId;
   
   glassArea = ((width - 140) / 1000.0) * (( height - 140) / 1000.0);
   
   RETURN COALESCE(glassArea, 0);
END;
$glassArea$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_frame_area(apertureId bigint)
RETURNS NUMERIC AS $frameArea$
declare
	area NUMERIC;
	glassArea NUMERIC;
	frameArea NUMERIC;
BEGIN
   SELECT aperture_area(apertureId), aperture_glass_area(apertureId)
   INTO area, glassArea 
   FROM apertures WHERE apertures.id = apertureId;
   
   frameArea = area - glassArea;
   
   RETURN COALESCE(frameArea, 0);
END;
$frameArea$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_overall_u_value(apertureId bigint)
RETURNS NUMERIC AS $overallUValue$
declare
	overallArea NUMERIC;
	glassArea NUMERIC;
	glassUValue NUMERIC;
	frameArea NUMERIC;
	frameUValue NUMERIC;
	overallUValue NUMERIC;
	customOverallUValue NUMERIC;
	isOverallUValueCustom BOOLEAN;
BEGIN
   SELECT COALESCE(apertures.overall_u_value, 0), apertures.is_overall_u_value_custom, aperture_area(apertureId), aperture_glass_area(apertureId), COALESCE(glasses.u_value, 0), aperture_frame_area(apertureId), COALESCE(profiles.u_value, 0)
   INTO customOverallUValue, isOverallUValueCustom, overallArea, glassArea, glassUValue, frameArea, frameUValue
   FROM apertures LEFT OUTER JOIN glasses ON apertures.glass_id = glasses.id LEFT OUTER JOIN profiles ON apertures.profile_id = profiles.id
   WHERE apertures.id = apertureId;
   
   IF isOverallUValueCustom THEN
   	RETURN customOverallUValue;
   END IF;
   
   IF overallArea = 0 THEN
   	RETURN 0; -- To prevent div by 0 error
   END IF;
   
   overallUValue = (glassUValue * glassArea + frameUValue * frameArea) / overallArea;
   
   RETURN COALESCE(overallUValue, 0);
END;
$overallUValue$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION aperture_total_discount(apertureId bigint)
RETURNS numeric AS $totalDiscount$
declare
	sellingPrice numeric;
	discountedSellingPrice numeric;
BEGIN
   	SELECT aperture_total_discounted_selling_price(apertureid), aperture_total_selling_price(apertureid)
	INTO discountedSellingPrice, sellingPrice
	FROM apertures
	WHERE apertures.id = apertureId;

	RETURN sellingPrice - discountedSellingPrice;
END;
$totalDiscount$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION quote_discount(quoteId BIGINT)
RETURNS numeric AS $discount$
declare
	customDiscount NUMERIC;
	isDiscountCustom BOOLEAN;
	sumDiscount NUMERIC;
BEGIN
	SELECT quotes.discount, quotes.is_discount_custom
	INTO customDiscount, isDiscountCustom
	FROM quotes
	WHERE quotes.id = quoteId;
	
	IF isDiscountCustom THEN
		RETURN customDiscount;
	END IF;

   	SELECT SUM(aperture_total_discount(apertures.id)) 
	INTO sumDiscount
	FROM quotes INNER JOIN apertures ON quotes.id = apertures.quote_id AND apertures.parent_aperture_id IS NULL AND apertures.deleted_at IS NULL
	WHERE quotes.id = quoteId;

	RETURN COALESCE(sumDiscount, 0);
END;
$discount$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION quote_gross_profit(quoteId BIGINT)
RETURNS numeric AS $grossProfit$
declare
	customGrossProfit NUMERIC;
	isGrossProfitCustom BOOLEAN;
	sumGrossProfit NUMERIC;
BEGIN
	SELECT quotes.gross_profit, quotes.is_gross_profit_custom
	INTO customGrossProfit, isGrossProfitCustom
	FROM quotes
	WHERE quotes.id = quoteId;
	
	IF isGrossProfitCustom THEN
		RETURN customGrossProfit;
	END IF;

   	SELECT SUM(aperture_total_profit(apertures.id)) 
	INTO sumGrossProfit
	FROM quotes INNER JOIN apertures ON quotes.id = apertures.quote_id AND apertures.parent_aperture_id IS NULL AND apertures.deleted_at IS NULL
	WHERE quotes.id = quoteId;

	RETURN COALESCE(sumGrossProfit, 0);
END;
$grossProfit$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION quote_sub_total(quoteId BIGINT)
RETURNS numeric AS $subTotal$
declare
	customSubTotal NUMERIC;
	isSubTotalCustom BOOLEAN;
	sumSubTotal NUMERIC;
BEGIN
	SELECT quotes.sub_total, quotes.is_sub_total_custom
	INTO customSubTotal, isSubTotalCustom
	FROM quotes
	WHERE quotes.id = quoteId;
	
	IF isSubTotalCustom THEN
		RETURN customSubTotal;
	END IF;

   	SELECT SUM(aperture_total_discounted_selling_price(apertures.id)) 
	INTO sumSubTotal
	FROM quotes INNER JOIN apertures ON quotes.id = apertures.quote_id AND apertures.parent_aperture_id IS NULL AND apertures.deleted_at IS NULL
	WHERE quotes.id = quoteId;

	RETURN COALESCE(sumSubTotal, 0);
END;
$subTotal$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION quote_vat(quoteId BIGINT)
RETURNS numeric AS $vat$
declare
	customVat NUMERIC;
	isVatCustom BOOLEAN;
	subTotal NUMERIC;
BEGIN
	SELECT quotes.vat, quotes.is_vat_custom, quote_sub_total(quoteId)
	INTO customVat, isVatCustom, subTotal
	FROM quotes
	WHERE quotes.id = quoteId;
	
	IF isVatCustom THEN
		RETURN customVat;
	END IF;

	RETURN COALESCE(subTotal * 0.18, 0);
END;
$vat$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION quote_total(quoteId BIGINT)
RETURNS numeric AS $total$
declare
	customTotal NUMERIC;
	isTotalCustom BOOLEAN;
	subTotal NUMERIC;
	vat NUMERIC;
BEGIN
	SELECT quotes.total, quotes.is_total_custom, quote_sub_total(quoteId), quote_vat(quoteId)
	INTO customTotal, isTotalCustom, subTotal, vat
	FROM quotes
	WHERE quotes.id = quoteId;
	
	IF isTotalCustom THEN
		RETURN customTotal;
	END IF;

	RETURN subTotal + vat;
END;
$total$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION quote_total_discounted_selling_price(quoteId BIGINT)
RETURNS numeric AS $totalDiscountedSellingPrice$
declare
	sumTotalDiscountedSellingPrice NUMERIC;
BEGIN
	SELECT SUM(aperture_total_discounted_selling_price(apertures.id))
	INTO sumTotalDiscountedSellingPrice
	FROM apertures
	WHERE apertures.quote_id = quoteId AND apertures.parent_aperture_id IS NULL AND apertures.deleted_at IS NULL;

	RETURN sumTotalDiscountedSellingPrice;
END;
$totalDiscountedSellingPrice$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION quote_total_area(quoteId BIGINT)
RETURNS numeric AS $totalArea$
declare
	sumTotalArea NUMERIC;
BEGIN
	SELECT SUM(aperture_total_area(apertures.id))
	INTO sumTotalArea
	FROM apertures
	WHERE apertures.quote_id = quoteId AND apertures.parent_aperture_id IS NULL AND apertures.deleted_at IS NULL;

	RETURN sumTotalArea;
END;
$totalArea$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION quote_updated_at(quoteId bigint)
RETURNS TIMESTAMP WITH TIME ZONE AS $updatedAt$
declare
	updatedAt TIMESTAMP WITH TIME ZONE;
BEGIN
	SELECT GREATEST(quotes.updated_at, (
		SELECT MAX(apertures.updated_at) as aperture_updated_at
		INTO updatedAt
		FROM apertures
		WHERE apertures.quote_id = quoteId AND apertures.deleted_at IS NULL
	))
	FROM quotes
	WHERE quotes.id = quoteId;

	RETURN updatedAt;
END;
$updatedAt$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION quote_updated_by(quoteId bigint)
RETURNS BIGINT AS $updatedBy$
declare
	updatedBy BIGINT;
BEGIN
   	SELECT apertures.updated_by
	INTO updatedBy
	FROM apertures
	WHERE apertures.updated_at = quote_updated_at(quoteId) AND (apertures.quote_id = quoteId AND apertures.deleted_at IS NULL)
	LIMIT 1;
	
	IF updatedBy IS NOT NULL THEN
		RETURN updatedBy;
	END IF;
	
	SELECT quotes.updated_by
	INTO updatedBy
	FROM quotes
	WHERE quotes.id = quoteId;

	RETURN updatedBy;
END;
$updatedBy$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION quote_nets_quantity(quoteId bigint)
RETURNS BIGINT AS $netsQuantity$
declare
	sumNetsQuantity BIGINT;
BEGIN
	SELECT SUM(aperture_nets_quantity(apertures.id) * apertures.quantity)
	INTO sumNetsQuantity
	FROM apertures
	WHERE apertures.quote_id = quoteId AND apertures.parent_aperture_id IS NULL AND apertures.deleted_at IS NULL;

	RETURN COALESCE(sumNetsQuantity, 0);
END;
$netsQuantity$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION quote_next_number()
RETURNS BIGINT AS $quoteNumber$
declare
	quoteNumber BIGINT;
BEGIN
   SELECT number 
   INTO quoteNumber 
   FROM quotes 
   WHERE quotes.number IS NOT NULL
   ORDER BY quotes.created_at DESC
   LIMIT 1;
   
   IF quoteNumber IS NULL THEN
		RETURN 1;
	END IF;
   
   RETURN quoteNumber + 1;
END;
$quoteNumber$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION quote_total_cost(quoteId BIGINT)
RETURNS numeric AS $totalCost$
declare
	sumTotalCost NUMERIC;
BEGIN
	SELECT SUM(aperture_total_cost(apertures.id))
	INTO sumTotalCost
	FROM apertures
	WHERE apertures.quote_id = quoteId AND apertures.parent_aperture_id IS NULL AND apertures.deleted_at IS NULL;

	RETURN sumTotalCost;
END;
$totalCost$ LANGUAGE plpgsql;